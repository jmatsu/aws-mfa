package main

import (
	"errors"
	"fmt"
	"github.com/jmatsu/aws-mfa/internal"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelWarn)

	codeRegexp := regexp.MustCompile("^\\d{6}$")

	app := &cli.App{
		Name:      "aws-mfa",
		Version:   fmt.Sprintf("%s (git revision %s)", internal.Version, internal.Commit),
		Copyright: "Jumpei Matsuda (@jmatsu)",
		Compiled:  internal.CompiledAt,
		Usage:     "Issue a new session with MFA devices.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "log-format",
				Action: func(_ *cli.Context, s string) error {
					switch s {
					case "json":
						logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
						slog.SetDefault(logger)
						return nil
					default:
						return fmt.Errorf("unknown format: %s", s)
					}
				},
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Action: func(_ *cli.Context, v string) error {
					level, ok := map[string]slog.Level{
						"debug": slog.LevelDebug,
						"info":  slog.LevelInfo,
						"warn":  slog.LevelWarn,
						"error": slog.LevelError,
					}[strings.ToLower(v)]

					if !ok {
						return fmt.Errorf("%s is not a valid for --log-level", v)
					}

					slog.SetLogLoggerLevel(level)
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "code",
				Aliases:  []string{"c"},
				Usage:    "Provide an OTP code.",
				Required: true,
				Action: func(_ *cli.Context, s string) error {
					if !codeRegexp.MatchString(s) {
						return fmt.Errorf("%s must be 6-digit number", s)
					}

					return nil
				},
			},
			&cli.IntFlag{
				Name:    "minutes",
				Aliases: []string{"m"},
				Usage:   "Specify how may minutes a new session works.",
				EnvVars: []string{"SESSION_TOKEN_MINUTES"},
				Value:   60,
				Action: func(_ *cli.Context, i int) error {
					if i < 15 {
						return fmt.Errorf("--minutes must be at least 15")
					} else if i > 2160 {
						return fmt.Errorf("--minutes must be at most 2160 (36 hours)")
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:     "mfa-profile",
				Usage:    "An aws profile that will have a new session.",
				EnvVars:  []string{"AWS_PROFILE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "without-mfa-profile",
				Usage: "An aws profile used to issue a new session.",
			},
			&cli.BoolFlag{
				Name:  "env",
				Usage: "Specify this if you would like to to use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY instead of specifying an aws profile.",
			},
			&cli.StringFlag{
				Name:    "device-arn",
				Usage:   "Specify the ARN of the device to issue the session to.",
				EnvVars: []string{"MFA_DEVICE"},
			},
		},
		Before: func(c *cli.Context) error {
			if !existAwsCli() {
				return fmt.Errorf("aws is not installed")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			input := &commandInput{
				Code:            c.String("code"),
				DurationSeconds: int32(c.Int("minutes")) * 60,
				PreferEnvVars:   c.Bool("env"),
			}

			if v := c.String("device-arn"); c.IsSet("device-arn") {
				input.DeviceArn = &v
			}

			input.ProfileWithMFA = c.String("mfa-profile")

			if v := c.String("without-mfa-profile"); c.IsSet("without-mfa-profile") {
				input.ProfileWithoutProfile = &v
			} else {
				v := fmt.Sprintf("%s-without-mfa", input.ProfileWithMFA)
				input.ProfileWithoutProfile = &v
			}

			if *input.ProfileWithoutProfile == "" {
				return fmt.Errorf("--without-mfa-profile is required")
			}

			return getAndSaveSessionToken(c.Context, input)
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())

		depth := 0

		for err != nil {
			slog.Info(strings.ReplaceAll(err.Error(), "\n", "。"), "err-depth", depth)

			err = errors.Unwrap(err)
		}

		os.Exit(1)
	}
}
