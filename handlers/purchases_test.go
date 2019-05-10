package handlers

import (
  // "fmt"
	"log"
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"
	"github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
)

func PerformTransactionCreation(account_id string, transaction_amount float64, operation_type int) *httptest.ResponseRecorder {
    gin.SetMode(gin.TestMode)
    handler := PerformPurchase
    router := gin.Default()
    router.POST("/v1/transactions", handler)

    newValidAccount := ReceivedPurchase{Account_id: account_id,
      Operation_type_id: operation_type,
      Amount: transaction_amount,
    }

    requestByte, _ := json.Marshal(newValidAccount)

    req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewReader(requestByte))

    if err != nil {
        log.Println(err)
    }

    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    return resp
}

func TestPurchaseTransactionInsuficientAmount(t *testing.T) {
  type IdOfValue struct {
    Id string `json:id`
  }

  var idValue IdOfValue
  accountCreated := PerformAccountCreation()
  err := json.Unmarshal(accountCreated.Body.Bytes(), &idValue)

  if (err != nil) {
    assert.Equal(t, false, true, "Could not create an account. Check your mongo connection.")
  }

  ourAccountId := idValue.Id

  resp := PerformTransactionCreation(ourAccountId, 9999999999, 2)

  assert.Equal(t, 400, resp.Code, "There should not be enough limit on your account")
}

func TestPurchaseTransactionWithInvalidOperationType(t *testing.T) {
  type IdOfValue struct {
    Id string `json:id`
  }

  var idValue IdOfValue
  accountCreated := PerformAccountCreation()
  err := json.Unmarshal(accountCreated.Body.Bytes(), &idValue)

  if (err != nil) {
    assert.Equal(t, false, true, "Could not create an account. Check your mongo connection.")
  }

  ourAccountId := idValue.Id

  resp := PerformTransactionCreation(ourAccountId, 10, 99)

  assert.Equal(t, 400, resp.Code, "This should not be a valid operation type value")
}


func TestPurchaseTransactionValidValue(t *testing.T) {
  type IdOfValue struct {
    Id string `json:id`
  }

  var idValue IdOfValue
  accountCreated := PerformAccountCreation()
  err := json.Unmarshal(accountCreated.Body.Bytes(), &idValue)

  if (err != nil) {
    assert.Equal(t, false, true, "Could not create an account. Check your mongo connection.")
  }

  ourAccountId := idValue.Id

  resp := PerformTransactionCreation(ourAccountId, 10, 2)

  assert.Equal(t, 200, resp.Code, "This should create a valid transaction")
}


func TestPurchaseTransactionValidWithdrawlValue(t *testing.T) {
  type IdOfValue struct {
    Id string `json:id`
  }

  var idValue IdOfValue
  accountCreated := PerformAccountCreation()
  err := json.Unmarshal(accountCreated.Body.Bytes(), &idValue)

  if (err != nil) {
    assert.Equal(t, false, true, "Could not create an account. Check your mongo connection.")
  }

  ourAccountId := idValue.Id

  withdrawlTypeTransaction := 3

  resp := PerformTransactionCreation(ourAccountId, 10, withdrawlTypeTransaction)

  assert.Equal(t, 200, resp.Code, "This should create a valid transaction")
}



func TestPurchaseTransactionInvalidAccountId(t *testing.T) {
  withdrawlTypeTransaction := 3
  invalidAccountId := "hamburguer"

  resp := PerformTransactionCreation(invalidAccountId, 10, withdrawlTypeTransaction)

  assert.Equal(t, 400, resp.Code, "This should return a bad request because it is an invalid id for account")
}
