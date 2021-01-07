#!/bin/sh

CMD_ARGS="-token $APP_TOKEN \
    -address 0.0.0.0 \
    -dataDir /srv/data \
    -mode $APP_MODE \
    -rateLimit $APP_RATE_LIMIT \
    -disableIPv6"

if [[ "$APP_TOKEN" == "" ]]; then
  echo "Token missing. Please provide one."
  exit 1
fi

if [[ "$APP_METRICS" != "" && "$APP_METRICS" != "false" ]]; then
  CMD_ARGS="$CMD_ARGS -metrics"
fi

./webhook2telegram $CMD_ARGS
