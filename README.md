# aws-mfa

A command to issue a new session token with mfa device authentication.

## Installation

*Recommended*

Download a binary from GitHub Release.

*Alternative*

```bash
go install github.com/jmatsu/aws-mfa@latest
```

## Usage

```
NAME:
   aws-mfa - Issue a new session with MFA devices.

USAGE:
   aws-mfa [global options] command [command options]

GLOBAL OPTIONS:
   --code value, -c value       Provide an OTP code.
   --minutes value, -m value    Specify how may minutes a new session works. (default: 60) [$SESSION_TOKEN_MINUTES]
   --mfa-profile value          An aws profile that will have a new session. [$AWS_PROFILE]
   --without-mfa-profile value  An aws profile used to issue a new session.
   --env                        Specify this if you would like to to use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY instead of specifying an aws profile. (default: false)
   --device-arn value           Specify the ARN of the device to issue the session to. [$MFA_DEVICE]
   --help, -h                   show help
   --version, -v                print the version
```

`aws-mfa` uses `--without-mfa-profile` to issue a new session token, and uses `--mfa-profile` value to store a new session token. You can call mfa-required API calls with the profile of `--mfa-profile` after running `aws-mfa`.

| --mfa-profile | --without-mfa-profile | A profile to issue a new session | A profile to call mfa-requested apis |
|:--------------|:----------------------|:---------------------------------|:-------------------------------------| 
| *unspecified* | *unspecified*         | N/a                              | N/a                                  |
| *unspecified* | bar                   | N/a                              | N/a                                  |
| foo           | *unspecified*         | foo                              | foo-without-mfa                      |
| foo           | bar                   | foo                              | bar                                  |

You don't have to specify `--device-arn` if you have only one MFA device. `aws-mfa` will find an ARN from `iam:ListMfaDevices` API. Please note that `--device-arn` is required if you have no MFA device or multiple MFA devices.

## Examples

```bash
aws-mfa --code <otp code> --mfa-profile foo
# or
AWS_PROFILE=foo aws-mfa --code <otp code>
```

- An access key to issue a new token must be stored in `foo-without-mfa` profile.
- A MFA device associated with your IAM will be chosen
- You can call MFA-required APIs by using `foo` profile's credentials.

```bash
aws-mfa --code <otp code>  --mfa-profile foo --env --device-arn arn:.....:mfa/...
```

- An access key to issue a new token must be retrieved from `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`. 
- `arn:.....:mfa/...` must associate with your IAM and `otp code`
- You can call MFA-required APIs by using `foo` profile's credentials.

# TODOs

- an assume-role arn to get an iam
- an assume-role arn to issue a session
- a session name when using an assume-role
