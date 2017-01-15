// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"fmt"
	"testing"

	"gopl.io/ch9/bank1"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("Deposit 200, Bank Balance=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		fmt.Println("Deposit 100, Bank Balance=", bank.Balance())
		done <- struct{}{}
	}()

	// John
	go func() {

		ok := bank.Withdraw(350)
		if ok {
			t.Errorf("Expected Withdraw of 350 to be %v", ok)
		}
		fmt.Println("Withdraw 350, Bank Balance=", bank.Balance())
		done <- struct{}{}
	}()

	// Jane
	go func() {
		ok := bank.Withdraw(50)
		if !ok {
			t.Errorf("Expected Withdraw of 50 to be %v", ok)
		}
		fmt.Println("Withdraw 50, Bank Balance=", bank.Balance())
		done <- struct{}{}
	}()

	// Wait for all transactions.
	<-done
	<-done
	<-done
	<-done

	if got, want := bank.Balance(), 250; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
	fmt.Println("Final Bank Balance=", bank.Balance())
}
