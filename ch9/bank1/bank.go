// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

var deposits = make(chan int)                      // send amount to deposit
var balances = make(chan int)                      // receive balance
var withdrawRequests = make(chan *WithdrawRequest) // send withdraw request

type WithdrawRequest struct {
	amount     int
	resultChan chan bool
}

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	wr := &WithdrawRequest{amount, make(chan bool)}
	// Send Request
	withdrawRequests <- wr
	return <-wr.resultChan
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case wr := <-withdrawRequests:
			if balance > wr.amount {
				balance -= wr.amount
				wr.resultChan <- true
			} else {
				wr.resultChan <- false
			}
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
