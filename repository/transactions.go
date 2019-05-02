package repository

import (

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

func FindDownedFrom(_id interface{}) []interface{} {
  mongoQuery := bson.D{{"related_purchase_transaction", _id.(primitive.ObjectID).Hex()}}

  results := findQuery(downedPurchaseEvent, mongoQuery)

  return _mapTo(results, "value")
}

func FindTransactionsNotIn(excluded []interface{}) []primitive.M {
  mongoQuery := bson.D{{"_id", bson.D{{"$nin", excluded}}}, {"operation.type", bson.D{{"$in", purchasesTypes}}}}

  results := findQuery(transactions, mongoQuery)

  return results //_mapTo(results, "_id")
}

func FindAllZeroedPurchaseTransactions() []interface{}{
  results := findAll(zeroedPurchaseEvent)

  parsed := _map(results, _toObjectIdFrom("transaction_id"))

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
