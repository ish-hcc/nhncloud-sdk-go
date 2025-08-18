package openstack

import (
	"fmt"
	"errors"
	"net/http"

	"github.com/cloud-barista/nhncloud-sdk-go"
)

// BaseError is an error type that all other error types embed.
type BaseError struct {
	DefaultErrString string
	Info             string
}

func (e BaseError) Error() string {
	e.DefaultErrString = "An error occurred while executing a Gophercloud request."
	return e.choseErrString()
}

func (e BaseError) choseErrString() string {
	if e.Info != "" {
		return e.Info
	}
	return e.DefaultErrString
}

// ErrEndpointNotFound is the error when no suitable endpoint can be found
// in the user's catalog
type ErrEndpointNotFound struct{ gophercloud.BaseError }

func (e ErrEndpointNotFound) Error() string {
	return "No suitable endpoint could be found in the service catalog."
}

// ErrInvalidAvailabilityProvided is the error when an invalid endpoint
// availability is provided
type ErrInvalidAvailabilityProvided struct{ gophercloud.ErrInvalidInput }

func (e ErrInvalidAvailabilityProvided) Error() string {
	return fmt.Sprintf("Unexpected availability in endpoint query: %s", e.Value)
}

// ErrNoAuthURL is the error when the OS_AUTH_URL environment variable is not
// found
type ErrNoAuthURL struct{ gophercloud.ErrInvalidInput }

func (e ErrNoAuthURL) Error() string {
	return "Environment variable OS_AUTH_URL needs to be set."
}

// ErrNoUsername is the error when the OS_USERNAME environment variable is not
// found
type ErrNoUsername struct{ gophercloud.ErrInvalidInput }

func (e ErrNoUsername) Error() string {
	return "Environment variable OS_USERNAME needs to be set."
}

// ErrNoPassword is the error when the OS_PASSWORD environment variable is not
// found
type ErrNoPassword struct{ gophercloud.ErrInvalidInput }

func (e ErrNoPassword) Error() string {
	return "Environment variable OS_PASSWORD needs to be set."
}

// ErrUnexpectedResponseCode is returned by the Request method when a response code other than
// those listed in OkCodes is encountered.
type ErrUnexpectedResponseCode struct {
	BaseError
	URL            string
	Method         string
	Expected       []int
	Actual         int
	Body           []byte
	ResponseHeader http.Header
}

// ResponseCodeIs returns true if this error is or contains an ErrUnexpectedResponseCode reporting
// that the request failed with the given response code. For example, this checks if a request
// failed because of a 404 error:
//
//	allServers, err := servers.List(client, servers.ListOpts{})
//	if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
//		handleNotFound()
//	}
//
// It is safe to pass a nil error, in which case this function always returns false.
func ResponseCodeIs(err error, status int) bool {
	var codeError ErrUnexpectedResponseCode
	if errors.As(err, &codeError) {
		return codeError.Actual == status
	}
	return false
}

// ErrTimeOut is the error type returned when an operations times out.
type ErrTimeOut struct {
	BaseError
}

func (e ErrTimeOut) Error() string {
	e.DefaultErrString = "A time out occurred"
	return e.choseErrString()
}
