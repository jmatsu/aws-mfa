#!/usr/bin/env bash

set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

usage() {
  cat <<EOF >&2
Usage:

  aws-mfa -h
  aws-mfa [-v] --code <code> [--minutes <minutes>] [--without-mfa-profile <aws profile>] [--mfa-profile <new aws profile>]

Issue a new session for a user who has enabled 2FA with a virtual device.

Options:
-h, --help             Print this help and exit
-v, --verbose          Print script debug info
-c, --code             Provide an OTP code.
-m, --minutes          Specify how may minues a new session works.
--mfa-profile          A aws profile that will have a new session.
--without-mfa-profile  A aws profile used to issue a new session.
EOF
  exit
}

cleanup() {
  trap - SIGINT SIGTERM ERR EXIT

  :
}

setup_colors() {
  if [[ -t 2 ]] && [[ -z "${NO_COLOR-}" ]] && [[ "${TERM-}" != "dumb" ]]; then
    # shellcheck disable=SC2034
    NOFORMAT='\033[0m' RED='\033[0;31m' GREEN='\033[0;32m' ORANGE='\033[0;33m' BLUE='\033[0;34m' PURPLE='\033[0;35m' CYAN='\033[0;36m' YELLOW='\033[1;33m'
  fi
}

msg() {
  echo >&2 -e "${1-}"
}

info() {
  msg "${GREEN-}$1${NOFORMAT-}"
}

warn() {
  msg "${YELLOW-}$1${NOFORMAT-}"
}

err() {
  msg "${RED-}$1${NOFORMAT-}"
}

die() {
  err "${1-}"
  exit "${2-1}"
}

parse_params() {
  minutes='60'
  code=''
  wo_mfa_profile=''
  mfa_profile="${AWS_PROFILE-}"

  _VERBOSE_=''

  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) _VERBOSE_=1 ;;
    --no-color) NO_COLOR=1 ;;
    -c | --code)
      code="${2-}"
      shift
      ;;
    -m | --minutes)
      minutes="${2-}"
      shift
      ;;
    --without-mfa-profile)
      wo_mfa_profile="${2-}"
      shift
      ;;
    --mfa-profile)
      mfa_profile="${2-}"
      shift
      ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac

    shift
  done

  if [[ -z "${mfa_profile:-}" ]]; then
    die "Missing required parameter: --mfa-profile or AWS_PROFILE env"
  fi

  if ((${minutes:-0} < 30)) || ((${minutes:-0} > 120)); then
    die "--minutes must be positive number and less than 120"
  fi

  if [[ -z "$code" ]]; then
    die "Missing required parameter: --code"
  fi

  if [[ ! "$code" =~ ^[0-9]{6}$ ]]; then
    die "OTP code must be 6-digits."
  fi

  return 0
}

parse_params "$@"
setup_colors

# shellcheck disable=SC2046
unset $(env | grep -E '^AWS_' | awk -F= '$0=$1')

if [[ -z "${wo_mfa_profile-}" ]]; then
  wo_mfa_profile="$mfa_profile-without-mfa"
fi

virtual_serial_arn=''
virtual_serial_arn="$(aws iam get-user --output json --profile "$wo_mfa_profile" | jq -r '.User.Arn' | sed -e 's/:user\//:mfa\//')"

if [[ -z "$virtual_serial_arn" ]]; then
  die "Cannot get the ARN of your virtual device."
fi

aws \
  sts \
  get-session-token \
  ${_VERBOSE_:+--debug} \
  --output json \
  --profile "$wo_mfa_profile" \
  --duration-seconds "$((minutes * 60))" \
  --serial-number "$virtual_serial_arn" \
  --token-code "$code" | \
  jq -r '"aws_access_key_id " + .Credentials.AccessKeyId, "aws_secret_access_key " + .Credentials.SecretAccessKey, "aws_session_token " + .Credentials.SessionToken, "expiration_date " + .Credentials.Expiration' | \
  xargs -n2 aws configure --profile "$mfa_profile" set
