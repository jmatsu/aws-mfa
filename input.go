package main

// commandInput is data collection of cli arguments
type commandInput struct {
	// Code is an OPT-code associated with DeviceArn
	Code string

	// DurationSeconds is a lifespan of a new token
	DurationSeconds int32

	// DeviceArn is an arn of a MFA device
	DeviceArn *string

	// ProfileWithoutProfile is a profile to issue a new token
	ProfileWithoutProfile *string

	// ProfileWithMFA is a profile to save a new token
	ProfileWithMFA string

	// PreferEnvVars is used to find credentials only from environment variables
	PreferEnvVars bool
}
