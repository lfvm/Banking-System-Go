package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)



func TestTransferTsx(t *testing.T){
	store := NewStore(testDb)


	account1:= createRandomAccount(t)
	account2:= createRandomAccount(t)
	fmt.Println(">> Before: " ,account1.Balance,account2.Balance )

	//run a concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTransactionResult)

	// Making a go routine for concurrent transactions
	for i := 0; i <n; i++{ 
		go func(){
			result, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountId: account1.ID,
				ToAccountId: account2.ID,
				Ammount: amount,
			})

			// sending errors and results to corresponding channels
			errs <- err
			results <- result
		}()
	}

	for i := 0; i <n; i++{ 
		err := <- errs
		require.NoError(t,err)
		
		result := <- results
		require.NotEmpty(t, result)


		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t,err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t,-amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t,amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)



		// check accounts 
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)		
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
	}

	// Check accounts balance 
	updateAccount1,err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t,err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)

	fmt.Println(">> After: " ,updateAccount1.Balance,updateAccount2.Balance )

	require.NoError(t,err)
	require.Equal(t, account1.Balance - int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n)*amount, updateAccount2.Balance)
}