package repository

import (
  // "log"
  "time"
)

type AccountCreditEvent struct {
	Account_id string
	Amount float64
	Transaction_id string
	Event_date time.Time
}

func CreateAccountCreditEvent(event AccountCreditEvent) (string, error){
  return insertOne(accountCreditEvent, event)
}
