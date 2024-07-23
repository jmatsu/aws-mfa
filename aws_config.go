package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"os"
)

var errCannotRetrieveAwsCredentials = fmt.Errorf("cannot retrieve AWS credentials")
var errCannotLoadDefaultConfig = fmt.Errorf("cannot load default config")
var errCannotLoadSharedConfig = fmt.Errorf("cannot load shared config")
var errCannotConfigureAccessKeyId = fmt.Errorf("cannot configure aws_access_key_id")
var errCannotConfigureSecretAccessKey = fmt.Errorf("cannot configure aws_secret_access_key")
var errCannotConfigureSessionToken = fmt.Errorf("cannot configure aws_session_token")
var errCannotConfigureExpirationDate = fmt.Errorf("cannot configure expiration_date")

type awsConfig struct {
	context.Context
	*aws.Config
}

// loadAwsConfig finds the config based on cli arguments a.k.a commandInput
func loadAwsConfig(ctx context.Context, input *commandInput) (*awsConfig, error) {
	var cfg aws.Config

	_ = os.Unsetenv("AWS_PROFILE")
	_ = os.Unsetenv("AWS_SESSION_TOKEN")

	if input.PreferEnvVars {
		if c, err := config.LoadDefaultConfig(ctx); err != nil {
			return nil, errors.Join(errCannotLoadDefaultConfig, err)
		} else {
			cfg = c
		}
	} else {
		_ = os.Unsetenv("AWS_ACCESS_KEY_ID")
		_ = os.Unsetenv("AWS_SECRET_ACCESS_KEY")

		if c, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(*input.ProfileWithoutProfile)); err != nil {
			return nil, errors.Join(errCannotLoadSharedConfig, err)
		} else {
			cfg = c
		}
	}

	if _, err := cfg.Credentials.Retrieve(ctx); err != nil {
		return nil, errors.Join(errCannotRetrieveAwsCredentials, err)
	}

	return &awsConfig{
		Context: ctx,
		Config:  &cfg,
	}, nil
}

// saveCredentials saves a given credentials via aws-cli calls
func (cfg *awsConfig) saveCredentials(profile string, cred *types.Credentials) error {
	cli := newAwsCli(cfg.Context)

	// Save to credentials
	if err := cli.ConfigureProp(profile, "aws_access_key_id", *cred.AccessKeyId); err != nil {
		return errors.Join(errCannotConfigureAccessKeyId, err)
	}

	if err := cli.ConfigureProp(profile, "aws_secret_access_key", *cred.SecretAccessKey); err != nil {
		return errors.Join(errCannotConfigureSecretAccessKey, err)
	}

	if err := cli.ConfigureProp(profile, "aws_session_token", *cred.SessionToken); err != nil {
		return errors.Join(errCannotConfigureSessionToken, err)
	}

	// Save to config
	if err := cli.ConfigureProp(profile, "expiration_date", cred.Expiration.String()); err != nil {
		return errors.Join(errCannotConfigureExpirationDate, err)
	}

	return nil
}
