package handlers

import (
	"log"
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

func PerformPayment(ctx *gin.Context) {
	alreadyZeroed := repository.FindAllZeroedPurchaseTransactions()
	notZeroed 		:= repository.FindTransactionsNotIn(alreadyZeroed)

	leftToPay := repository.Map(notZeroed, func(value bson.M) interface{} {
		total := value["amount"].(float64)
		paid  := utils.ReduceNumbers(repository.FindDownedFrom(value["_id"]), utils.Sum)
		return TransactionIdAndAmountLeft{AmountLeft: total-paid, Id: value["_id"].(primitive.ObjectID)}
	})

	log.Printf("%v", leftToPay)

	ctx.JSON(200, gin.H{
		"still_active": notZeroed,
	})
}
