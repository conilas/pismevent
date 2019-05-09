package repository

import (
  // "log"
  "time"
  "go.mongodb.org/mongo-driver/bson"
)

type AccountWithdrawlEvent struct {
	Account_id string
	Amount float64
	Transaction_id string
	Event_date time.Time
}

func FindWithdrawlEventsAmountsFromAccount(_id string) []interface{} {
  mongoQuery := bson.D{{"account_id", _id}}

  results := findQuery(accountWthdrawlEvent, mongoQuery)

  return _mapTo(results, "amount")
}

func CreateAccountWithdrawlEvent(event AccountWithdrawlEvent) (string, error){
  return insertOne(accountWthdrawlEvent, event)
}
