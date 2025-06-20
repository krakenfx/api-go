package helper

// Must returns the provided value if err is nil; otherwise, it panics with the error.
// Useful for writing tests to reduce error handling boilerplate.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
