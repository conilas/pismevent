package repository

import (
  "log"
  "go.mongodb.org/mongo-driver/bson"
)

type Account struct {
  Id string

  Available_credit_limit float64

  Available_withdraw_limit float64
}

func FindAllAccounts() []bson.M{
  return findAll(accounts)
}

func FindAccountFromId(_id string)  (Account, error){

  var acc Account

  err := findOneById(accounts, _id).Decode(&acc)

  if err != nil {
    log.Printf("[ERROR] Could not parse", err)
    return  Account{}, err
  }

  return Account{Available_withdraw_limit: acc.Available_withdraw_limit,
    Available_credit_limit: acc.Available_credit_limit,
    Id: _id,
  }, nil
 
}
