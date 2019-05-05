package repository

import (
  "time"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Operation struct {
  Type string `json: type`

  Charge_order float64
}

type Transaction struct {
  _id primitive.ObjectID

  Account_id string

  Operation Operation

  Amount float64

  Event_date time.Time

  Due_date time.Time
}

func CreateTransaction(transaction Transaction) (string, error)  {
  return insertOne(transactions, transaction)
}

func CreatePaymentTransaction(transaction Transaction) (string, error)  {
  transaction.Operation = Operation{Type: "PAYMENT", Charge_order: 0}

  return insertOne(transactions, transaction)
}

func FindUnzeroedPayments(excluded []interface{}) []primitive.M {
  mongoQuery := bson.D{{"_id", bson.D{{"$nin", excluded}}}, {"operation.type", bson.D{{"$in", paymentTypes}}}}

  results := findQuery(transactions, mongoQuery)

  sortedResults := sortByUrgency(results)

  return sortByUrgency(sortedResults)
}

func FindUnpaidTransactions(excluded []interface{}) []primitive.M {
  mongoQuery := bson.D{{"_id", bson.D{{"$nin", excluded}}}, {"operation.type", bson.D{{"$in", purchasesTypes}}}}

  results := findQuery(transactions, mongoQuery)

  sortedResults := sortByUrgency(results)

  return sortByUrgency(sortedResults)
}

func FindAllTransactions() []interface{}{
  results := _mapTo(findAll(transactions), "amount" )

  return results
}
