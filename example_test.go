package retry_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/azr/retry"
)

var (
	ErrSomething = errors.New("Holly guacamole, something failed.")
)

func TestCall(t *testing.T) {
	nb := &NoopBackoff{}
	beenThere := 0
	maxRetry := 3
	err := retry.Call{
		Fn: func() error {
			beenThere++
			return ErrSomething
		},
		OnRetry:       nb.Backoff,
		MaxRetry:      maxRetry,
		IsRetryableFn: func(_ error) bool { return true },
	}.Run()
	if err != ErrSomething {
		t.Errorf("Wrong error returned, expected %v, got %v", ErrSomething, err)
	}
	if beenThere != maxRetry+1 {
		t.Errorf("Did not retry enough !! retryed : %d, expected %d", beenThere, maxRetry+1)
	}
}

func TestCallNoIsRetryableFn(t *testing.T) {
	nb := &NoopBackoff{}
	beenThere := 0
	maxRetry := 3
	err := retry.Call{
		Fn: func() error {
			beenThere++
			return ErrSomething
		},
		OnRetry:  nb.Backoff,
		MaxRetry: maxRetry,
	}.Run()
	if err != ErrSomething {
		t.Errorf("Wrong error returned, expected %v, got %v", ErrSomething, err)
	}
	if beenThere != maxRetry+1 {
		t.Errorf("Did not retry enough !! retryed : %d, expected %d", beenThere, maxRetry+1)
	}
}

func TestCallNoRetry(t *testing.T) {
	nb := &NoopBackoff{}
	beenThere := 0
	maxRetry := 3
	err := retry.Call{
		Fn: func() error {
			beenThere++
			return ErrSomething
		},
		OnRetry:       nb.Backoff,
		MaxRetry:      maxRetry,
		IsRetryableFn: func(_ error) bool { return false },
	}.Run()
	if err != ErrSomething {
		t.Errorf("Wrong error returned, expected %v, got %v", ErrSomething, err)
	}
	if beenThere != 1 {
		t.Errorf("Did not retry enough !! retryed : %d, expected %d", beenThere, 1)
	}
}

func ExampleCall() {
	beenThere := 0
	err := retry.Call{
		Fn: func() error {
			fmt.Println("Been there:", beenThere)
			if beenThere == 2 {
				return nil // third call is successfull !
			}
			beenThere++
			return ErrSomething
		},
		OnRetry: func() {
			fmt.Println("Retrying...")
			// sleep some time
		},
		MaxRetry:      2,
		IsRetryableFn: nil, // meaning everything is retryable
	}.Run()
	fmt.Println(err)
	// Output: Been there: 0
	// Retrying...
	// Been there: 1
	// Retrying...
	// Been there: 2
	// <nil>
}

type NoopBackoff struct {
	called int
}

func (n *NoopBackoff) Backoff() { n.called++ }
