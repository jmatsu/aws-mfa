package main

import (
	"context"
	"os/exec"
)

func existAwsCli() bool {
	_, err := exec.LookPath("aws")
	return err == nil
}

type awsCli interface {
	ConfigureProp(profile string, name string, value string) error
}

func newAwsCli(ctx context.Context) awsCli {
	return &defaultAwsCli{
		Context: ctx,
	}
}

type defaultAwsCli struct {
	context.Context
}

func (cli *defaultAwsCli) ConfigureProp(profile string, name string, value string) error {
	cmd := exec.CommandContext(cli.Context, "aws", "configure", "--profile", profile, "set", name, value)
	return cmd.Run()
}
