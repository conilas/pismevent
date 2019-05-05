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
	Id     string  `json:"id"`

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
		moneyLeft 				  = moneyLeft - eventAmount

		downedEvent := repository.DownedEvent{Related_purchase_transaction: currentTransaction.Id.Hex(),
			Related_payment_transaction: paymentTransactionId, Value: eventAmount,
			Event_date: time.Now(),
		}
		paymentDownedId, _  := repository.CreateDownedPaymentEvent(downedEvent)
		purchaseDownedId, _  := repository.CreateDownedPurchaseEvent(downedEvent)

		if (moneyLeft == 0) {
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

func processPayment(payment ReceivedPayment) float64{
	_   					    = repository.FindAccountFromId(payment.Id)
	alreadyZeroed    := repository.FindAllZeroedPurchaseTransactions()
	notZeroed 		   := repository.FindUnpaidTransactions(alreadyZeroed)

	transactionId, _ := repository.CreatePaymentTransaction(repository.Transaction{Account_id: payment.Id,
		Amount: payment.Amount, Event_date: time.Now()})

	unpaidTransactions := repository.Map(notZeroed, func(value bson.M) interface{} {
		total := value["amount"].(float64)
		alreadyPaid  := utils.ReduceNumbers(repository.FindDownedFrom(value["_id"]), utils.Sum)
		return TransactionIdAndAmountLeft{AmountLeft: total-alreadyPaid, Id: value["_id"].(primitive.ObjectID)}
	})

	return processPurchaseDownings(payment.Amount, unpaidTransactions, transactionId)
}

func PerformPayment(ctx *gin.Context) {

  receivedPayments := parsePayment(ctx)
	var moneyLeft float64

	for _, receivedPayment := range receivedPayments {
		log.Printf("[DEBUG] Received: %v", receivedPayment)
		moneyLeft = processPayment(receivedPayment)
		if (moneyLeft == 0 ) { break }
	}

	log.Printf("[DEBUG] Heya! After all payments, we still have %v left.", moneyLeft)

	//TODO create an event with putting the moneyLeft on the creditLimit for the person

	ctx.JSON(200, gin.H{
		"still_active": "",
	})
}
