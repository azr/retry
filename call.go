// Small pkg that helps you retry things that can be.
//
//  Given a func Fn that might fail and return an error.
//   if the error `IsRetryAble`, Fn will be rerun.
//   until MaxRetry is reached.
//
// If you had a foo lib for a foo service you could:
//   package foo
//
//   import (
//       "github.com/azr/backoff"
//       "github.com/azr/retry"
//   )
//
//   func Retry(fn func() error) retry.Call {
//     return retry.Call{
//         Fn:            fn,
//         OnRetry:       GetBackoff().BackOff,
//         IsRetryableFn: IsErrorRetryable,
//         MaxRetry:      2,
//     }
//   }
//
//   func GetBackoff() backoff.Interface {
//      return backoff.NewExponential()
//       configure backoffer
//   }
//
//   func IsErrorRetryable(err error) bool {
//       if err == SomeRetryableError {
//           return true
//       }
//       return false
//   }
//
//   func SomeFunc() error {
//       err := log.Output(1, "foo")
//       return err
//   }
//
// Now you can just :
//
//   package main
//
//   import "bar/foo"
//
//   func main() {
//       err := foo.Retry(foo.SomeFunc).Run()
//   }
//
package retry

type Call struct {
	Fn            func() error     // Called upon Run
	IsRetryableFn func(error) bool // Called when Fn returns an error, with the error. Default behavior make this always true
	OnRetry       func()           // Called if there was an error (retryable) and if defined
	MaxRetry      int              // Fn can be called up to `1 + MaxRetry` retryable times
}

func (r Call) Run() error {
	err := r.Fn()
	if err != nil && r.MaxRetry > 0 {
		if r.IsRetryableFn == nil {
			return r.rerun()
		} else if r.IsRetryableFn(err) {
			return r.rerun()
		}
	}
	return err
}

func (r Call) rerun() error {
	r.MaxRetry--
	if r.OnRetry != nil {
		r.OnRetry()
	}
	return r.Run()
}
