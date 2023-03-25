package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yankycranky/my-bank/util"
)

func CreateMyAccount(t *testing.T, owner *string) Account {
	var currOwner string
	if owner != nil {

		currOwner = *owner
	} else {
		currOwner = util.RandomOwner()
	}
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    currOwner,
		Balance:  util.RandomMoney(),
		Currency: "USD",
	})
	if err != nil {
		log.Fatal("Cannot Create Account", err)
	}
	assert.NoError(t, err)
	return account
}

func TestCreateAccount(t *testing.T) {

	account := CreateMyAccount(t, nil)
	assert.Equal(t, account.Currency, "USD")
}

func TestGetAccount(t *testing.T) {
	str := "Ankit"
	account := CreateMyAccount(t, &str)
	assert.Equal(t, account.Currency, "USD")
	account, err := testQueries.GetAccountByName(context.Background(), str)
	if err != nil {
		log.Fatal("Cannot retrieve account")
	}
	t.Log(account)

	assert.Equal(t, account.Owner, "Ankit")
	// assert.Equal(t, account.Balance, "20000")
	assert.Equal(t, account.Currency, "USD")
}

func TestTransferAccount(t *testing.T) {

	acc1 := CreateMyAccount(t, nil)
	acc2 := CreateMyAccount(t, nil)
	fmt.Printf("current Balances 1=> %v and 2=> %v \n\n", acc1.Balance, acc2.Balance)
	myStore := NewStore(dbConn)
	errs := make(chan error)
	results := make(chan TransferResult)
	n := 5
	for i := 0; i < n; i++ {
		go func(i int) {
			result, err := myStore.TransferTx(context.Background(), TransferRequest{
				from_account:       acc1.ID,
				to_account:         acc2.ID,
				from_account_owner: acc1.Owner,
				to_account_owner:   acc2.Owner,
				amount:             200,
			})

			errs <- err
			results <- result
		}(i)
	}
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		fmt.Printf("Final Result >>> 1=>> %v and 2=>>%v\n", result.FromAccount.Balance, result.ToAccount.Balance)
		assert.NoError(t, err)
		// val, _ := strconv.Atoi(acc2.Balance)
		assert.Equal(t, result.success, true)
		assert.Equal(t, result.ToEntry.Amount, strconv.Itoa(200))
		assert.Equal(t, result.FromEntry.Amount, strconv.Itoa(-200))
		assert.Equal(t, result.Transfer.FromAccountID, acc1.ID)
		assert.Equal(t, result.FromEntry.AccountID, acc1.ID)
		assert.Equal(t, result.Transfer.ToAccountID, acc2.ID)
		assert.Equal(t, result.ToEntry.AccountID, acc2.ID)

		// Check Accounts
		assert.NotEmpty(t, result.FromAccount)
		assert.NotEmpty(t, result.ToAccount)
		oldBal1, _ := strconv.Atoi(acc1.Balance)
		assert.Equal(t, strconv.Itoa(oldBal1-200*(i+1)), result.FromAccount.Balance)
		oldBal2, _ := strconv.Atoi(acc2.Balance)
		assert.Equal(t, strconv.Itoa(oldBal2+200*(i+1)), result.ToAccount.Balance)
		fmt.Printf("\n\n")
	}
}
