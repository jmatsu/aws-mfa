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
```

todo:

Allow to specify

- a device arn
- an assume-role arn to get an iam
- an assume-role arn to issue a session
- a session name when using an assume-role