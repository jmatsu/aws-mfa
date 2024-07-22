# aws-mfa

A aws CLi wrapper to issue a new session for a user who has enabled 2FA.

*Non-production ready.*

- Currently, this works only when no assume-role is required and a mfa device is a first virtual device

## Usage

```
  aws-mfa [-v] --code <code> [--minutes <minutes>] [--without-mfa-profile <aws profile>] [--mfa-profile <new aws profile>]

Options:
-h, --help             Print this help and exit
-v, --verbose          Print script debug info
-c, --code             Provide an OTP code.
-m, --minutes          Specify how may minues a new session works.
--mfa-profile          A aws profile that will have a new session.
--without-mfa-profile  A aws profile used to issue a new session.
--env                  Specify this if you would like to to use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY instead of specifying an aws profile.
```

--mfa-profile | --without-mfa-profile | A profile to issue a new session
:--- |:---|:--- 
`<null>` | `<null>` | N/A
"foo" | `<null>`| "foo-without-mfa"
"foo" | "bar" | "bar"
`<null>` | "bar" | N/A

`--mfa-profile` is required but the value of `--mfa-profile` defaults to `AWS_PROFILE` environment variable unless specified.

AWS_PROFILE | --mfa-profile | A profile for mfa-required APIs
:--- |:---|:--- 
`<null>` | `<null>` | N/A
"foo" | `<null>` | "foo"
"foo" | "bar" | "bar"
`<null>` | "bar" | "bar"

todo:

Allow to specify

- a device arn
- an assume-role arn to get an iam
- an assume-role arn to issue a session
- a session name when using an assume-role
