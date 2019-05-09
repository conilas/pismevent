# Pismevent


## Idea

The idea is to perform the challenge but in a event-sourced way. In order to do so, *every* event will be stored in the database - including amounts payed, drawbacked transactions and so on. The application is simple, but as it is for payments and bank related stuff, there should be a tracking of everything - and I do mean everything.

The approach would also be simple if it wasn't for the validations one must do for each of those transactions. I'll show some of them below and explain how they were overcome.

## Example 1: Incoming payment (installment or in cash)

Whenever we have an incoming payment, the applications must check if there's a credit transaction that wasn't yet zeroed and, if so, create some related events. Let's say that, in this case, we do have some values that were already credited and we must validate it. In order to do so, it must:

* Find which was not yet zeroed: go find every transaction with a specific type that indicates payment that aren't yet in **zeroed_purchase_event** collection. So it performs a *not in* operation.
* If it finds any that wasn't yet zeroed, it must check how much is still has left. In order to do so, it goes in the **downed_purchase_event** collection and sees how much was taken from that transaction.
* From that value, it will check how much it must discount from the incoming transaction. That means it will insert a new value in the **downed_purchase_event** and then create the other events related to the transaction.

## Example 2: Incoming purchase 

Whenever we have an incoming purchase, we will follow a flow that is almost like the payment one, except this is a bit easier because payment *may* have more than one payment per request **and** it must order to see which transactions were killed. So, in a way, we:

* First, check if the operation type is credit-related or withdrawl-related. Create an event on the necessary collection (either **account_withdrawl_event** or **account_credit_event**).
* Find which still have some credit: go find every transaction with a specific type that indicates payment that aren't yet in **downed_payment_event** collection. So it performs a *not in* operation.
* If it finds any that wasn't yet zeroed, it must check how much is still has left. In order to do so, it goes in the **zeroed_payment_event** collection and sees how much was taken from that transaction.
* From that value, it will check how much it must discount from the incoming transaction. That means it will insert a new value in the **zeroed_payment_event** and then create the other events related to the transaction.

### Not good enough? Interested? Check the doc, which contains two flowcharts to explain the idea.

  
  
