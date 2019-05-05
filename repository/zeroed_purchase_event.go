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

func CreateZeroedPurchaseEvent(event ZeroedPurchaseEvent) (string, error){
  return insertOne(zeroedPurchaseEvent, event)
}
