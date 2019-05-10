package handlers

import (
  "fmt"
	"log"
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"
  "pismevent/repository"
	"github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
)

//treat this as a suit of tests - it will test the whole flux related to account
//creation, listing, patching & asserting the values are as expected
func TestAccountCreationAndPatchingItAssertingWithExpectedValue(t *testing.T){
  type IdOfValue struct {
    Id string `json:id`
  }

  initialCredit      := float64(9999)  //hardcoded just 4 testing
  initialWithdrawl   := float64(9999)  //hardcoded just 4 testing
  decreasedCredit    := float64(-2000) //hardcoded just 4 testing
  decreasedWithdrawl := float64(-2000) //hardcoded just 4 testing


  var idValue IdOfValue
  accountCreated := PerformAccountCreation()
  err := json.Unmarshal(accountCreated.Body.Bytes(), &idValue)

  if (err != nil) {
    assert.Equal(t, false, true, "Could not create an account. Check your mongo connection.")
  }

  ourAccountId := idValue.Id

  resp := SendPatchRequest(ourAccountId, decreasedCredit, decreasedWithdrawl)

  assert.Equal(t, 200, resp.Code)

  log.Printf("[DEBUG] Our id: ", ourAccountId)

  allAccountsNow := PerformAccountListing()

  type AllAccounts struct {
    Accounts []repository.Account `json:accounts`
  }

  var accounts AllAccounts

  err = json.Unmarshal(allAccountsNow.Body.Bytes(), &accounts)

  if (err != nil) {
    assert.Equal(t, false, true, "Could not list all accounts. Check your mongo connection.")
  }

  for _, account := range accounts.Accounts {
    if (account.Id == ourAccountId) {
      assert.Equal(t, account.Available_credit_limit, decreasedWithdrawl + initialWithdrawl,
        fmt.Sprintf("Withdrawl should be decreased and now equal to %v", decreasedWithdrawl + initialWithdrawl))

      assert.Equal(t, account.Available_credit_limit, decreasedCredit + initialCredit,
        fmt.Sprintf("Credit should be decreased and now equal to %v", decreasedCredit + initialCredit))
    }
  }
}

func PerformAccountListing() *httptest.ResponseRecorder {
  log.Printf("[DEBUG] Remeber: this test will throw error if there are no mongo instance. Please, provide one.")

  gin.SetMode(gin.TestMode)
  handler := RetrieveAccountsLimits
  router := gin.Default()
  router.GET("/v1/accounts", handler)

  req, err := http.NewRequest("GET", "/v1/accounts", nil)

  if err != nil {
      log.Println(err)
  }

  resp := httptest.NewRecorder()

  router.ServeHTTP(resp, req)

  return resp
}

func PerformAccountCreation() *httptest.ResponseRecorder{
  log.Printf("[DEBUG] Remeber: this test will throw error if there are no mongo instance. Please, provide one.")

  gin.SetMode(gin.TestMode)
  handler   := CreateAccount
  router    := gin.Default()
  router.POST("/v1/accounts", handler)

  newValidAccount := LimitChanger{Available_credit_limit: Amount{Amount: 9999},
                                  Available_withdrawl_limit: Amount{Amount: 9999}}

  requestByte, _ := json.Marshal(newValidAccount)

  req, err := http.NewRequest("POST", "/v1/accounts", bytes.NewReader(requestByte))

  if err != nil {
      log.Println(err)
  }

  resp := httptest.NewRecorder()
  router.ServeHTTP(resp, req)

  return resp
}

func SendPatchRequest(id string, credit_limit_amount float64, withdrawl_limit_amount float64) *httptest.ResponseRecorder{
  log.Printf("[DEBUG] Remeber: this test will throw error if there are no mongo instance. Please, provide one.")

  gin.SetMode(gin.TestMode)
  handler   := PerformActionOnAccount
  router    := gin.Default()
  router.PATCH("/v1/accounts/:id", handler)

  newValidAccount := LimitChanger{Available_credit_limit: Amount{Amount: credit_limit_amount},
                                  Available_withdrawl_limit: Amount{Amount: withdrawl_limit_amount}}

  requestByte, _ := json.Marshal(newValidAccount)

  req, err := http.NewRequest("PATCH", "/v1/accounts/" + id, bytes.NewReader(requestByte))

  if err != nil {
      log.Println(err)
  }

  resp := httptest.NewRecorder()

  router.ServeHTTP(resp, req)

  return resp
}

func TestAccountCreation(t *testing.T) {
  resp := PerformAccountCreation()
  assert.Equal(t, 200, resp.Code)
}

func TestPatchAccountLimitWithInvalidId(t *testing.T){
  resp := SendPatchRequest("hamburguer", 9999, 9999)
  assert.Equal(t, 400, resp.Code)
}

func TestListAccountsByListingAllAccounts(t *testing.T) {
    resp := PerformAccountListing()
    assert.Equal(t, resp.Code, 200)
}
