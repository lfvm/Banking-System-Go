package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTransaction(ctx context.Context, fn func(*Queries) error) error { 

	// This function is used to execute transactions in a db 
	// a transaction in a db is when there are required multiple operations 
	// in the db. Since one of this operations can fail. we would want to go back to 
	// the original state of the db. 
	// if everyting is okay, then we commit the changes to the db 

	tx, err := store.db.BeginTx(ctx,nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err :%v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()

}	

type TransferTransactionParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId int64 `json:"to_account_id"`
	Ammount int64 `json:"ammount"`
}

type TransferTransactionResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

// performs a money transfer from one account to other 
// It creates a transfer record, add account entries
// and update account within a single db transaction 
func (store *SQLStore) TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error) {

	var result TransferTransactionResult
	err := store.execTransaction(ctx, func(q *Queries) error {
		
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID: arg.ToAccountId,
			Amount: arg.Ammount,
		})	
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: -arg.Ammount,
		})
	
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount: arg.Ammount,
		})

		if err != nil {
			return err
		}

		// To preven deadlock from happening on the database, we can make sure to always update the 
		// account with the smaller id first. Lets remember that the order in which queries happen is important 
		// and can make deadlock errors if not done properly 
		if arg.FromAccountId < arg.ToAccountId { 
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Ammount, arg.ToAccountId,arg.Ammount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Ammount, arg.FromAccountId, -arg.Ammount)
			if err != nil {
				return err
			}
		}

		
		return nil
	})

	return result, err
}


func addMoney(
	ctx context.Context, 
	q *Queries,
	accountId1 int64, 
	amount int64, 
	accountId2 int64,
	ammount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId1,
		Ammount: amount,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId2,
		Ammount: ammount2,
	})

	if err != nil {
		return
	}
	return 
}
