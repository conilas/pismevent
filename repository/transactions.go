package repository

import (
  "log"
  "time"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Operation struct {
  _type string `json: type`

  charge_order float64
}

type Transaction struct {
  _id primitive.ObjectID

  account_id string

  operation Operation

  amount float64

  event_date time.Time

  due_date time.Time
}

func FindDownedFrom(_id interface{}) []interface{} {
  mongoQuery := bson.D{{"related_purchase_transaction", _id.(primitive.ObjectID).Hex()}}

  results := findQuery(downedPurchaseEvent, mongoQuery)

  return _mapTo(results, "value")
}

func FindTransactionsNotIn(excluded []interface{}) []primitive.M {
  mongoQuery := bson.D{{"_id", bson.D{{"$nin", excluded}}}, {"operation.type", bson.D{{"$in", purchasesTypes}}}}

  results := findQuery(transactions, mongoQuery)

  sortedResults := sortByUrgency(results)

  log.Printf("%v", sortedResults)

  log.Printf("%v", _mapTo(sortedResults, "_id"))

  return sortByUrgency(sortedResults) 
}

func FindAllZeroedPurchaseTransactions() []interface{}{
  results := findAll(zeroedPurchaseEvent)

  parsed := Map(results, _toObjectIdFrom("transaction_id"))

  return parsed
}

func FindAllZeroedPaymentTransactions() []interface{}{
  results := _mapTo(findAll(zeroedPaymentEvent), "_id" )

  return results
}

func FindAllTransactions() []interface{}{
  results := _mapTo(findAll(transactions), "amount" )

  return results
}
