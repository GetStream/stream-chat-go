package stream

import "errors"

var (
	// ErrorMissingChannelType indicates that the function requires a valid
	// channel type argument.
	ErrorMissingChannelType = errors.New("channel type is empty")

	// ErrorMissingUserID indicates that the function requires a valid
	// user ID argument.
	ErrorMissingUserID = errors.New("user ID is empty")

	// ErrorMissingTargetID indicates that the function requires a valid
	// target ID argument.
	ErrorMissingTargetID = errors.New("target ID is empty")

	// ErrorEmptyTargetID indicates that the function requires one or more
	// target IDs.
	ErrorEmptyTargetID = errors.New("target ID list is empty")

	// ErrorMissingDeviceID indicates that the function requires a valid device ID.
	ErrorMissingDeviceID = errors.New("device ID is empty")
)
