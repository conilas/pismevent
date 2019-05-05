package main

import (
	handlers "eventsourcismo/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/v1/payments", handlers.PerformPayment)
	r.Run(":3031") // listen and serve on 0.0.0.0:8080
}
