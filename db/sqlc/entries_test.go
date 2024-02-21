package db

import (
	"context"
	"testing"

	"github.com/lfvm/simplebank/utils"
	"github.com/stretchr/testify/require"
)


func createRandomEntry(t *testing.T) Entry{

	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: utils.RandomMoney(),
	}

	entry,err := testQueries.CreateEntry(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,entry)
	require.Equal(t,arg.AccountID,entry.AccountID)
	require.Equal(t,arg.Amount,entry.Amount)
	require.NotZero(t,entry.ID)
	require.NotZero(t,entry.CreatedAt)
	return entry
}

func deleteEntryWithId(t *testing.T, id int64){
	err:= testQueries.DeleteEntry(context.Background(),id)
	require.NoError(t,err)

	entry, err := testQueries.GetEntry(context.Background(),id)
	require.Empty(t, entry)
	require.EqualError(t,err, "sql: no rows in result set")
}

func TestCreateEntry(t *testing.T){
	entry := createRandomEntry(t)
	deleteAccountWithId(t, entry.ID)
}

func TestDeleteEntry(t *testing.T){
	entry := createRandomEntry(t)
	deleteEntryWithId(t, entry.ID)
}

func TestListEntries(t *testing.T){

	createdEntries := make([]Entry,10)

	for i:=0; i<10; i++{
		createdEntries[i] = createRandomEntry(t)
	}
	arg := ListEntriesParams{
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,entries,5)
	for _,entry := range entries{
		require.NotEmpty(t,entry)
	}
	for _,entry := range createdEntries{
		deleteAccountWithId(t, entry.ID)
		deleteEntryWithId(t, entry.ID)
	}
}

func TestGetEntryById(t *testing.T){
	entry := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(),entry.ID)
	require.NoError(t,err)
	require.NotEmpty(t,entry2)
	require.Equal(t,entry,entry2)
	deleteAccountWithId(t, entry.ID)
	deleteEntryWithId(t, entry.ID)
}