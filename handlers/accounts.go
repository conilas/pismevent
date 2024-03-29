package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
	repository "pismevent/repository"
  utils "pismevent/utils"
)

type Account struct {
  Id string

  Available_credit_limit float64

  Available_withdraw_limit float64
}

type Amount struct {
	Amount float64 `json:"amount"`
}

type LimitChanger struct {
	Available_credit_limit Amount `json:"available_credit_limit"`
	Available_withdrawl_limit Amount `json:"available_withdrawl_limit"`
}

func parseLimitChanger(ctx *gin.Context) LimitChanger {
	var limits LimitChanger
  ctx.BindJSON(&limits)
  return limits
}

func CreateAccount(ctx *gin.Context) {
	account   := parseLimitChanger(ctx)
	accountId, err := repository.CreateAccount(repository.AccountCreation{Available_credit_limit: account.Available_credit_limit.Amount,
		Available_withdraw_limit: account.Available_withdrawl_limit.Amount})

	if (err != nil){
		ctx.JSON(500, gin.H{
			"message": err,
		})

		return
	}

	ctx.JSON(200, gin.H{
		"id": accountId,
	})
}

func PerformActionOnAccount(ctx *gin.Context) {
	accountId      := ctx.Param("id")
	_, err   := repository.FindAccountFromId(accountId)

	if (err != nil ) {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Invalid account id: %v", accountId),
		})

		return
	}

	limitChanger := parseLimitChanger(ctx)

	repository.CreateAccountCreditEvent(repository.AccountCreditEvent{
		Account_id: accountId, Amount: limitChanger.Available_credit_limit.Amount, Event_date: time.Now(),
	})

	repository.CreateAccountWithdrawlEvent(repository.AccountWithdrawlEvent{
		Account_id: accountId, Amount: limitChanger.Available_withdrawl_limit.Amount, Event_date: time.Now(),
	})

	ctx.JSON(200, gin.H{
		"accounts": ctx.Param("id"),
	})
}

func RetrieveAccountsLimits(ctx *gin.Context) {
  accounts             := repository.FindAllAccounts()
  currentStateAccounts := repository.Map(accounts, func(account bson.M) interface{} {
    accountId                   := account["_id"].(primitive.ObjectID).Hex()
    allEventsFromAccount        := repository.FindEventsAmountsFromAccount(accountId)
		allWithdrawlEvents 				  := repository.FindWithdrawlEventsAmountsFromAccount(accountId)
		valueInTotal 				        := utils.ReduceNumbers(allEventsFromAccount, utils.Sum)
		withdrawlValueInTotal 			:= utils.ReduceNumbers(allWithdrawlEvents, utils.Sum)
    currentAccountCreditLimit   := account["available_credit_limit"].(float64) + valueInTotal
    currentAccountWithdrawLimit := account["available_withdraw_limit"].(float64) + withdrawlValueInTotal

    return Account{Id: accountId,
			Available_credit_limit: currentAccountCreditLimit,
			Available_withdraw_limit: currentAccountWithdrawLimit,
		}
  })

  ctx.JSON(200, gin.H{
		"accounts": currentStateAccounts,
	})
}
