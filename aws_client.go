package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
)

var errCannotListMyMFADevices = fmt.Errorf("failed to list mfa devices from iam service (iam:ListVirtualMFADevices)")
var errCannotGetSessionToken = fmt.Errorf("failed to get a session token from sts service (sts:GetSessionToken)")

type awsClient interface {
	ListMyMFADevice() ([]iamTypes.MFADevice, error)
	GetSessionToken(request *sts.GetSessionTokenInput) (*stsTypes.Credentials, error)
}

type defaultAwsClient struct {
	context.Context
	iamClient *iam.Client
	stsClient *sts.Client
}

func newAwsClient(ctx context.Context, cfg *awsConfig) (awsClient, error) {
	return &defaultAwsClient{
		Context:   ctx,
		iamClient: iam.NewFromConfig(*cfg.Config),
		stsClient: sts.NewFromConfig(*cfg.Config),
	}, nil
}

func (client *defaultAwsClient) ListMyMFADevice() ([]iamTypes.MFADevice, error) {
	out, err := client.iamClient.ListMFADevices(client.Context, &iam.ListMFADevicesInput{})

	if err != nil {
		return nil, errors.Join(errCannotListMyMFADevices, err)
	}

	return out.MFADevices, nil
}
func (client *defaultAwsClient) GetSessionToken(request *sts.GetSessionTokenInput) (*stsTypes.Credentials, error) {
	out, err := client.stsClient.GetSessionToken(client.Context, request)

	if err != nil {
		return nil, errors.Join(errCannotGetSessionToken, err)
	}

	return out.Credentials, nil
}
