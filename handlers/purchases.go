package handlers

import (
	"fmt"
	"log"
	"math"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
	"eventsourcismo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	repository "eventsourcismo/repository"
)

var ValidationPerOperation = map[int] func(repository.Account, float64) (float64, error){
  1: validatePurchaseAmount,
  2: validatePurchaseAmount,
  3: validateWithdrawlAmount,
}

var ProcessPerOperation = map[int] func(repository.Transaction) (float64){
  1: processPurchase,
  2: processPurchase,
  3: processWithdrawl,
}

type ReceivedPurchase struct {
	Account_id     string     `json:"account_id"`
  Operation_type_id int `json:"Operation_type_id"`
	Amount float64    `json:"amount"`
}

func parsePurchase(ctx *gin.Context) ReceivedPurchase {
	var purchase ReceivedPurchase
  ctx.BindJSON(&purchase)
  return purchase
}

func processPaymentDownings(initial float64, unpaidTransactions []interface{},purchaseTransactionId string) float64{
	moneyLeft := initial

	for _, stillCreditedTransaction := range unpaidTransactions {
		currentTransaction := stillCreditedTransaction.(TransactionIdAndAmountLeft)
		eventAmount 			 := math.Min(currentTransaction.AmountLeft, moneyLeft)
    perfectFit         := currentTransaction.AmountLeft == moneyLeft
		moneyLeft 				  = moneyLeft - eventAmount

		downedEvent := repository.DownedEvent{Related_purchase_transaction: purchaseTransactionId,
			Related_payment_transaction: currentTransaction.Id.Hex(), Value: eventAmount,
			Event_date: time.Now(),
		}
		paymentDownedId, _  := repository.CreateDownedPaymentEvent(downedEvent)
		purchaseDownedId, _  := repository.CreateDownedPurchaseEvent(downedEvent)

		if (moneyLeft == 0) {

			if (perfectFit) {
        repository.CreateZeroedPaymentEvent(repository.ZeroedPaymentEvent{Transaction_id: currentTransaction.Id.Hex(),
  				Event_date: time.Now(), Last_downed_event: paymentDownedId})
      }

      repository.CreateZeroedPurchaseEvent(repository.ZeroedPurchaseEvent{Transaction_id: purchaseTransactionId,
				Event_date: time.Now(), Last_downed_event: purchaseDownedId})
			break
		} else {
      repository.CreateZeroedPaymentEvent(repository.ZeroedPaymentEvent{Transaction_id: currentTransaction.Id.Hex(),
				Event_date: time.Now(), Last_downed_event: paymentDownedId})
		}
	}

	return moneyLeft
}

func processWithdrawl(payment repository.Transaction) (float64){

	transactionId, _ := repository.CreateTransaction(payment)

	repository.CreateAccountWithdrawlEvent(repository.AccountWithdrawlEvent{Account_id: payment.Account_id,
		Amount: - payment.Amount, Transaction_id: transactionId, Event_date: time.Now(),
	})

	return 0
}

func processPurchase(payment repository.Transaction) (float64){
  alreadyZeroed    := repository.FindAllZeroedPaymentTransactions()
	notZeroed 		   := repository.FindUnzeroedPayments(alreadyZeroed)

	transactionId, _ := repository.CreateTransaction(payment)

	transactionsWithCredit := repository.Map(notZeroed, func(value bson.M) interface{} {
		total := value["amount"].(float64)
		alreadyPaid  := utils.ReduceNumbers(repository.FindDownedPaymentsFrom(value["_id"]), utils.Sum)
		return TransactionIdAndAmountLeft{AmountLeft: total-alreadyPaid, Id: value["_id"].(primitive.ObjectID)}
	})

	repository.CreateAccountCreditEvent(repository.AccountCreditEvent{Account_id: payment.Account_id,
		Amount: - payment.Amount, Transaction_id: transactionId, Event_date: time.Now(),
	})

	return processPaymentDownings(payment.Amount, transactionsWithCredit, transactionId)
}

func validateWithdrawlAmount(account repository.Account, needed float64) (float64, error) {
	allWithdrawlEvents 				  := repository.FindWithdrawlEventsAmountsFromAccount(account.Id)
	withdrawlValueInTotal 			:= utils.ReduceNumbers(allWithdrawlEvents, utils.Sum)
	currentAccountWithdrawLimit := account.Available_withdraw_limit + withdrawlValueInTotal

	if (currentAccountWithdrawLimit < needed) {
		return 0, errors.New("Insuficient amount")
	}

	return currentAccountWithdrawLimit, nil
}

func validatePurchaseAmount(account repository.Account, needed float64) (float64, error){
	allEventsFromAccount      := repository.FindEventsAmountsFromAccount(account.Id)
	valueInTotal 				      := utils.ReduceNumbers(allEventsFromAccount, utils.Sum)
	currentAccountCreditLimit := account.Available_credit_limit + valueInTotal

	if (currentAccountCreditLimit < needed) {
		return 0, errors.New("Insuficient amount")
	}

	return currentAccountCreditLimit, nil
}

func PerformPurchase(ctx *gin.Context) {

  receivedPurchase := parsePurchase(ctx)

	account, err   			 := repository.FindAccountFromId(receivedPurchase.Account_id)

	if (err != nil) {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Invalid account id: %v", receivedPurchase.Account_id),
		})

		return
	}

	//TODO create operation type validation

	currentAccountCreditLimit, err := ValidationPerOperation[receivedPurchase.Operation_type_id](account, receivedPurchase.Amount)

	if (err != nil) {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Not enough limit on account. You have: %v", currentAccountCreditLimit),
		})

		return
	}

	moneyLeft        := ProcessPerOperation[receivedPurchase.Operation_type_id] (repository.Transaction{Account_id: receivedPurchase.Account_id,
    Operation: utils.AllowedTransactionTypes[receivedPurchase.Operation_type_id],
    Amount: receivedPurchase.Amount, Event_date: time.Now()})


	log.Printf("[DEBUG] Heya! After all payments, we still have %v left.", moneyLeft)

	ctx.JSON(200, gin.H{
		"current_account_limit": currentAccountCreditLimit - receivedPurchase.Amount,
	})
}
