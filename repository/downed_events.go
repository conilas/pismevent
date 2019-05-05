package repository

import (
  "time"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)


type DownedEvent struct {
  Related_purchase_transaction string
  Related_payment_transaction  string
  Value float64
  Event_date time.Time
}

func FindDownedPaymentsFrom(_id interface{}) []interface{} {
  mongoQuery := bson.D{{"related_payment_transaction", _id.(primitive.ObjectID).Hex()}}

  results := findQuery(downedPaymentEvent, mongoQuery)

  return _mapTo(results, "value")
}

func FindDownedPurchasesFrom(_id interface{}) []interface{} {
  mongoQuery := bson.D{{"related_purchase_transaction", _id.(primitive.ObjectID).Hex()}}

  results := findQuery(downedPurchaseEvent, mongoQuery)

  return _mapTo(results, "value")
}

func CreateDownedPurchaseEvent(event DownedEvent) (string, error){
  return insertOne(downedPurchaseEvent, event)
}

func CreateDownedPaymentEvent(event DownedEvent) (string, error){
  return insertOne(downedPaymentEvent, event)
}
