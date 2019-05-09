package utils

import (
  "eventsourcismo/repository"
)

var AllowedTransactionTypes = map[int] repository.Operation{
  1: repository.Operation{Type: "IN_CASH_PURCHASE", Charge_order: 2},
  2: repository.Operation{Type: "INSTALLMENT_PURCHASE", Charge_order: 1},
  3: repository.Operation{Type: "WITHDRAWL", Charge_order: 0},
}
