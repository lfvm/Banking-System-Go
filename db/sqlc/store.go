package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTransaction(ctx context.Context, fn func(*Queries) error) error { 

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
	ToEntry int64 `json:"to_entry"`
}

// performs a money transfer from one account to other 
// It creates a transfer record, add account entries
// and update account within a single db transaction 


// func (store *Store) TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, err) {

	
// }