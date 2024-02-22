package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lfvm/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTrasnfer(t *testing.T) Transfer{

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: utils.RandomMoney(),
	}	

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	
	require.NoError(t,err)
	require.NotEmpty(t,transfer)
	require.Equal(t,arg.FromAccountID,transfer.FromAccountID)
	require.Equal(t,arg.Amount,transfer.Amount)
	require.Equal(t,arg.ToAccountID,transfer.ToAccountID)
	require.NotZero(t,transfer.ID)
	require.NotZero(t,transfer.CreatedAt)
	return transfer
}

func deleteTranserById(t *testing.T, id int64){
	err:= testQueries.DeleteTransfer(context.Background(),id)
	require.NoError(t,err)

	transfer, err := testQueries.GetTransfer(context.Background(),id)
	require.EqualError(t,err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)

}


func TestCreateTransfer(t *testing.T){
	transfer:= createRandomTrasnfer(t)
	deleteTranserById(t, transfer.ID)
	deleteAccountWithId(t, transfer.FromAccountID)
	deleteAccountWithId(t, transfer.ToAccountID)
}

func TestDeleteTransfer(t *testing.T){
	transfer := createRandomTrasnfer(t)
	deleteTranserById(t,transfer.ID)
	deleteAccountWithId(t, transfer.FromAccountID)
	deleteAccountWithId(t, transfer.ToAccountID)
}

func TestListTransfers(t *testing.T){

	createdTransfers  := make([]Transfer, 10)
	
	for i:=0; i<10; i++{
		createdTransfers[i] = createRandomTrasnfer(t)
	}

	arg := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,transfers,5)

	for _, transfer := range createdTransfers{
		require.NotEmpty(t,transfer)
		deleteTranserById(t, transfer.ID)
		deleteAccountWithId(t, transfer.FromAccountID)
		deleteAccountWithId(t, transfer.ToAccountID)
	}
}