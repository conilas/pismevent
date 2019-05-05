package repository

import (
  // "log"
  "time"
)

type ZeroedPaymentEvent struct {
  Transaction_id string
  Event_date time.Time
  Last_downed_event  string
}

func CreateZeroedPaymentEvent(event ZeroedPaymentEvent) (string, error){
  return insertOne(zeroedPaymentEvent, event)
}
