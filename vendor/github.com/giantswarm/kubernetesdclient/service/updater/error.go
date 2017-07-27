package updater

import "github.com/juju/errgo"

var executionFailedError = errgo.New("execution failed")

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return errgo.Cause(err) == executionFailedError
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var notFoundError = errgo.New("not found")

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return errgo.Cause(err) == notFoundError
}
