package readiness

type Option func(*Readiness)

// With pull failed handler func.
func WithPullFailedHandler(handler func(string, error)) Option {
	return func(r *Readiness) {
		r.onPullFailed = handler
	}
}
