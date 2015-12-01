// Small pkg that helps you retry things that can be.
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
