package handlers

import (
	"log"
	"math"
	"time"
	"github.com/gin-gonic/gin"
	"eventsourcismo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	repository "eventsourcismo/repository"
)

type TransactionIdAndAmountLeft struct {
	Id primitive.ObjectID

	AmountLeft float64
}

type ReceivedPayments []ReceivedPayment

type ReceivedPayment struct {
	Account_id     string  `json:"account_id"`

	Amount float64 `json:"amount"`
}

func parsePayment(ctx *gin.Context) ReceivedPayments {
	var payment ReceivedPayments
  ctx.BindJSON(&payment)
  return payment
}

func processPurchaseDownings(initial float64, unpaidTransactions []interface{},paymentTransactionId string) float64{
	moneyLeft := initial

	for _, unpaidTransactions := range unpaidTransactions {
		currentTransaction := unpaidTransactions.(TransactionIdAndAmountLeft)
		eventAmount 			 := math.Min(currentTransaction.AmountLeft, moneyLeft)
		perfectFit         := currentTransaction.AmountLeft == moneyLeft
		moneyLeft 				  = moneyLeft - eventAmount

		downedEvent := repository.DownedEvent{Related_purchase_transaction: currentTransaction.Id.Hex(),
			Related_payment_transaction: paymentTransactionId, Value: eventAmount,
			Event_date: time.Now(),
		}
		paymentDownedId, _  := repository.CreateDownedPaymentEvent(downedEvent)
		purchaseDownedId, _  := repository.CreateDownedPurchaseEvent(downedEvent)

		if (moneyLeft == 0) {

			if (perfectFit) {
				repository.CreateZeroedPurchaseEvent(repository.ZeroedPurchaseEvent{Transaction_id: currentTransaction.Id.Hex(),
					Event_date: time.Now(), Last_downed_event: purchaseDownedId})
			}

			repository.CreateZeroedPaymentEvent(repository.ZeroedPaymentEvent{Transaction_id: paymentTransactionId,
				Event_date: time.Now(), Last_downed_event: paymentDownedId})
			break
		} else {
			repository.CreateZeroedPurchaseEvent(repository.ZeroedPurchaseEvent{Transaction_id: currentTransaction.Id.Hex(),
				Event_date: time.Now(), Last_downed_event: purchaseDownedId})
		}
	}

	return moneyLeft
}

func processPayment(payment ReceivedPayment) (float64, string) {
	_   					    = repository.FindAccountFromId(payment.Account_id)
	//TODO: validate account existance
	alreadyZeroed    := repository.FindAllZeroedPurchaseTransactions()
	notZeroed 		   := repository.FindUnpaidTransactions(alreadyZeroed)

	transactionId, _ := repository.CreatePaymentTransaction(repository.Transaction{Account_id: payment.Account_id,
		Amount: payment.Amount, Event_date: time.Now()})

	unpaidTransactions := repository.Map(notZeroed, func(value bson.M) interface{} {
		total := value["amount"].(float64)
		alreadyPaid  := utils.ReduceNumbers(repository.FindDownedPurchasesFrom(value["_id"]), utils.Sum)
		return TransactionIdAndAmountLeft{AmountLeft: total-alreadyPaid, Id: value["_id"].(primitive.ObjectID)}
	})

	return processPurchaseDownings(payment.Amount, unpaidTransactions, transactionId), transactionId
}


func PerformPayment(ctx *gin.Context) {

  receivedPayments := parsePayment(ctx)

	for _, receivedPayment := range receivedPayments {
		log.Printf("[DEBUG] Received: %v", receivedPayment)
		_, transactionId := processPayment(receivedPayment)

		repository.CreateAccountCreditEvent(repository.AccountCreditEvent{Account_id: receivedPayment.Account_id,
			Amount: receivedPayment.Amount, Transaction_id: transactionId, Event_date: time.Now(),
		})
	}

	ctx.JSON(200, gin.H{
		"still_active": "",
	})
}
