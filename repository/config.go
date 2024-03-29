package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var client *mongo.Client
var ctx context.Context
var db *mongo.Database
var transactions *mongo.Collection
var zeroedPaymentEvent *mongo.Collection
var zeroedPurchaseEvent *mongo.Collection
var downedPaymentEvent *mongo.Collection
var downedPurchaseEvent *mongo.Collection
var accounts *mongo.Collection
var accountCreditEvent *mongo.Collection
var accountWthdrawlEvent *mongo.Collection

func mountAllCollections() {
	transactions = db.Collection("transactions")
	zeroedPaymentEvent = db.Collection("zeroed_payment_event")
	zeroedPurchaseEvent = db.Collection("zeroed_purchase_event")
	downedPaymentEvent = db.Collection("downed_payment_event")
	downedPurchaseEvent = db.Collection("downed_purchase_event")
	accounts = db.Collection("accounts")
	accountCreditEvent = db.Collection("account_credit_event")
	accountWthdrawlEvent = db.Collection("account_withdrawl_event")
}

func init() {
	_client, _ := mongo.NewClient(options.Client().ApplyURI(mongoIp))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = _client.Connect(ctx)
	db = _client.Database("pismo")
	client = _client
	mountAllCollections()
}
