package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/lfvm/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomInt(0, 1000),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func deleteAccountWithId(t *testing.T, id int64) {

	err := testQueries.DeleteAccount(context.Background(), id)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), id)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestCreateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	deleteAccountWithId(t, account1.ID)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	deleteAccountWithId(t, account1.ID)
}

func TestUpdateAccount(t *testing.T) {

	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: account1.Balance + 1000,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance+1000, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	deleteAccountWithId(t, account1.ID)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	deleteAccountWithId(t, account1.ID)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	var accountsCreated [10]Account

	for i := 0; i < 10; i++ {
		accountsCreated[i] = createRandomAccount(t)
		lastAccount = accountsCreated[i]
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

	for _, account := range accountsCreated {
		deleteAccountWithId(t, account.ID)
	}

}
