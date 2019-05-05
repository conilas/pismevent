package repository

import (
  // "log"
  "time"
)

type ZeroedPurchaseEvent struct {
  Transaction_id string
  Event_date time.Time
  Last_downed_event  string
}


func FindAllZeroedPurchaseTransactions() []interface{}{
  results := findAll(zeroedPurchaseEvent)

  parsed := Map(results, _toObjectIdFrom("transaction_id"))

  return parsed
}

func CreateZeroedPurchaseEvent(event ZeroedPurchaseEvent) (string, error){
  return insertOne(zeroedPurchaseEvent, event)
}
