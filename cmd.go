package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var errNoMFADeviceFound = fmt.Errorf("no mfa device was found")
var errMultipleMFADeviceFound = fmt.Errorf("only one mfa device is allowed")

// getAndSaveSessionToken is a core logic
func getAndSaveSessionToken(ctx context.Context, input *commandInput) error {
	cfg, err := loadAwsConfig(ctx, input)

	if err != nil {
		return err
	}

	aws, err := newAwsClient(ctx, cfg)

	if err != nil {
		return err
	}

	if input.DeviceArn == nil {
		if devices, err := aws.ListMyMFADevice(); err != nil {
			return err
		} else if len(devices) == 0 {
			return errNoMFADeviceFound
		} else if len(devices) > 1 {
			return errMultipleMFADeviceFound
		} else {
			input.DeviceArn = devices[0].SerialNumber
		}
	}

	cred, err := aws.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber:    input.DeviceArn,
		TokenCode:       &input.Code,
		DurationSeconds: &input.DurationSeconds,
	})

	if err != nil {
		return err
	}

	return cfg.saveCredentials(input.ProfileWithMFA, cred)
}
