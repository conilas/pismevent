package main

import (
	handlers "pismevent/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/v1/payments", handlers.PerformPayment)
	r.POST("/v1/transactions", handlers.PerformPurchase)
	r.POST("/v1/accounts", handlers.CreateAccount)
	r.GET("/v1/accounts/limits", handlers.RetrieveAccountsLimits)
	r.PATCH("/v1/accounts/:id", handlers.PerformActionOnAccount)
	r.Run(":3031") // listen and serve on 0.0.0.0:8080
}
