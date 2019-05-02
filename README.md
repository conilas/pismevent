# Pismevent


## Idea

The idea is to perform the challenge but in a event-sourced way. In order to do so, *every* event will be stored in the database - including amounts payed, drawbacked transactions and so on. The application is simple, but as it is for payments and bank related stuff, there should be a tracking of everything - and I do mean everything.

The approach would also be simple if it wasn't for the validations one must do for each of those transactions. I'll show some of them below and explain how they were overcome.

## Example 1: Incoming payment (installment or in cash)

Whenever we have an incoming payment, the applications must check if there's a credit transaction that wasn't yet zeroed and, if so, create some related events. Let's say that, in this case, we do have some values that were already credited and we must validate it. In order to do so, it must:

* Find which was not yet zeroed: go find every transaction with a specific type that indicates payment that aren't yet in **already_zeroed_transaction** collection. So it performs a *not in* operation.
* If it finds any that wasn't yet zeroed, it must check how much is still has left. In order to do so, it goes in the **transaction_credit_discounted** collection and sees how much was taken from that transaction.
* From that value, it will check how much it must discount from the incoming transaction. That means it will insert a new value in the **transaction_credit_discounted** and then create the other events related to the transaction.


Como achar transações não abatidas: Todas que estão em transactions mas que não estão em already_drawbacked_transactions.
  Como achar seu valor para abatimento: Montar o histórico vindo da transaction_drawbacks;
  Quando é necessário achar seu valor abatido: Quando entra uma nova transação de crédito;
  
Como achar transações com valor a ser abatido: Todas que estão em transactions mas não estão em already_zeroed_transactions.
  Como achar seu valor de crédito: Montar o histórico vindo da transaction_credit;
  Quando é necessário achar seu valor de crédito: Quando entra uma nova transação de pagamento;   
  
  
