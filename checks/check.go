package checks

// Check is a health/readiness check.
type Check func() error
