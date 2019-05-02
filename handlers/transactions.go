package handlers

import (
	"log"
	"github.com/gin-gonic/gin"
	"eventsourcismo/utils"
	repository "eventsourcismo/repository"
)

func FindUndownedPurchases(ctx *gin.Context) {
	alreadyZeroed := repository.FindAllZeroedPurchaseTransactions()
	notZeroed 		:= repository.FindTransactionsNotIn(alreadyZeroed)

	for _, value := range notZeroed {
		total := value["amount"].(float64)
		paid  := utils.ReduceNumbers(repository.FindDownedFrom(value["_id"]), utils.Sum)

		log.Printf("[DEBUG] Total: %v", total)
		log.Printf("[DEBUG] Already paid: %v", paid)
		log.Printf("[DEBUG] Missing: %v", (total - paid))
	}

	ctx.JSON(200, gin.H{
		"still_active": notZeroed,
	})
}
