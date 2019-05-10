## Event-sourced applications

The idea here was to develop a program that could control user transactions - such as payments and purchase - and keep track of the  limits of the refered account. 

So, as the challenge specifies, for every new transaction, we should keep an internal value that says whether or not it was credited/still had some debt left. The thing is - in the proposal, it says to update the referenced transaction whenever we have some new one - say, decrease the amount left to pay or something. But if we did that, we'd have:

1. Mutable state, which is hard to keep track of;
2. No way (or maybe a hard and forced way) to check what was on the db on a specified date;

So, in order to avoid that, I tried an event sourced approach, in which:

1. Everything that happens is an event: when there is a new payment transaction, the application goes in the DB, finds a non-yet payed purchase transaction and creates:
  * An event for the transaction;
  * An event for the account credit limit;
  * An event for every downings it does to that transaction;
  * An event for every downings in the current transaction;
2. No updates, meaning that we can have every history mounted at any determined point in time: IMO this is pretty usefull for a bank - you could, for instance, easily show the operations of the user in the last month and to which transaction each one was related;
3. ~~The current state is shown using Kafka or some event-streaming library; (did not have enough time to implement this)~~

This means we will have many collections on the DB, which are:

* ```account_credit_event``` - representing any event related to the account (patch on the account credit limit, for instance) in terms of credit;
* ```account_withdrawl_event``` - same as above, but only for withdrawl; 
* ```accounts``` - which keeps only the initial state of the account (when it was created);
* ```downed_payment_event``` and ```downed_purchase_event``` - every downing on a purchase/payment transaction coming from the other instance (i.e a payment that subtracts some of a purchase debt will be stored on the downed_payment_event);
* ```transaction``` - will keep track of every direct transaction created;
* ```zeroed_payment_event``` and ```zeroed_purchase_event``` - will keep track of every transaction that was zeroed (i.e a purchase that now has no debit on the DB).

Some sources:

* https://www.martinfowler.com/eaaDev/EventSourcing.html

* The Clean Architecture book on event sourcing
