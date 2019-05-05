package repository

import (
  // "log"
  "time"
  "go.mongodb.org/mongo-driver/bson"
)

type AccountCreditEvent struct {
	Account_id string
	Amount float64
	Transaction_id string
	Event_date time.Time
}

func FindEventsAmountsFromAccount(_id string) []interface{} {
  mongoQuery := bson.D{{"account_id", _id}}

  results := findQuery(accountCreditEvent, mongoQuery)

  return _mapTo(results, "amount")
}

func CreateAccountCreditEvent(event AccountCreditEvent) (string, error){
  return insertOne(accountCreditEvent, event)
}
