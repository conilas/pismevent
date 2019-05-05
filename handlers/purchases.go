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

func processPurchase(payment repository.Transaction) float64{
	_   					    = repository.FindAccountFromId(payment.Account_id)
  //TODO: validate account existance
	alreadyZeroed    := repository.FindAllZeroedPaymentTransactions()
	notZeroed 		   := repository.FindUnzeroedPayments(alreadyZeroed)

	transactionId, _ := repository.CreateTransaction(payment)

	transactionsWithCredit := repository.Map(notZeroed, func(value bson.M) interface{} {
		total := value["amount"].(float64)
		alreadyPaid  := utils.ReduceNumbers(repository.FindDownedPaymentsFrom(value["_id"]), utils.Sum)
		return TransactionIdAndAmountLeft{AmountLeft: total-alreadyPaid, Id: value["_id"].(primitive.ObjectID)}
	})

	return processPaymentDownings(payment.Amount, transactionsWithCredit, transactionId)
}

func PerformPurchase(ctx *gin.Context) {

  receivedPurchase := parsePurchase(ctx)

  //TODO ask what is the due_date for
	moneyLeft        := processPurchase(repository.Transaction{Account_id: receivedPurchase.Account_id,
    Operation: utils.AllowedTransactionTypes[receivedPurchase.Operation_type_id],
    Amount: receivedPurchase.Amount, Event_date: time.Now()})

	log.Printf("[DEBUG] Heya! After all payments, we still have %v left.", moneyLeft)

	//TODO create an event with putting the moneyLeft on the creditLimit for the person

	ctx.JSON(200, gin.H{
		"still_active": "",
	})
}
