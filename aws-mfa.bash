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
-m, --minutes          Specify how may minues a new session works. (SESSION_TOKEN_MINUTES)
--mfa-profile          A aws profile that will have a new session. (AWS_PROFILE)
--without-mfa-profile  A aws profile used to issue a new session.
--env                  Specify this if you would like to to use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY instead of specifying an aws profile.
--device-arn           An ARN of a virtual/phisical MFA device. (MFA_DEVICE)
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
  minutes="${SESSION_TOKEN_MINUTES-60}"
  code=''
  wo_mfa_profile=''
  mfa_profile="${AWS_PROFILE-}"
  read_env=''
  device_arn="${MFA_DEVICE-}"

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
    --env) read_env=1 ;;
    --device-arn)
      device_arn="${2-}"
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

  if [[ -n "$read_env" ]] && [[ -n "${wo_mfa_profile}" ]]; then
    die "Conflicted parameters: --env and --without-mfa-profile cannot be specified together"
  fi

  if [[ -n "$read_env" ]]; then
    if [[ -z "${AWS_ACCESS_KEY_ID-}" ]]; then
      die "Missing environment variable: AWS_ACCESS_KEY_ID is required when --env is specified."
    fi

    if [[ -z "${AWS_SECRET_ACCESS_KEY-}" ]]; then
      die "Missing environment variable: AWS_SECRET_ACCESS_KEY is required when --env is specified."
    fi
  fi

  if (($minutes < 15)) || (($minutes > 2160)); then
    die "--minutes must be from 15 mins to 2160 mins (36 hours)"
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

# We don't want cli to evaluate these variables.
unset AWS_PROFILE AWS_SESSION_TOKEN

common_options=()

if [[ -n "${_VERBOSE_-}" ]]; then
  common_options+=(--debug)
fi

mfa_options=(--profile "$mfa_profile")
without_mfa_options=()

if [[ -z "$read_env" ]]; then
  if [[ -z "${wo_mfa_profile-}" ]]; then
    without_mfa_options+=("--profile" "$mfa_profile-without-mfa")
  else
    without_mfa_options+=("--profile" "$wo_mfa_profile")
  fi
fi

if [[ -z "$device_arn" ]]; then
  device_arn="$(aws iam get-user --output json "${without_mfa_options[@]}" | jq -r '.User.Arn' | sed -e 's/:user\//:mfa\//')"
fi

if [[ -z "$device_arn" ]]; then
  die "Cannot get the ARN of your virtual device."
fi

aws \
  sts \
  get-session-token \
  "${common_options[@]}" \
  "${without_mfa_options[@]}" \
  --output json \
  --duration-seconds "$((minutes * 60))" \
  --serial-number "$device_arn" \
  --token-code "$code" | \
  jq -r '"aws_access_key_id " + .Credentials.AccessKeyId, "aws_secret_access_key " + .Credentials.SecretAccessKey, "aws_session_token " + .Credentials.SessionToken, "expiration_date " + .Credentials.Expiration' | \
    xargs -n2 aws configure \
      "${common_options[@]}" \
      "${mfa_options[@]}" \
      set
