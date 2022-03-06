#!/bin/sh

CMD_ARGS="-address 0.0.0.0 \
    -dataDir /srv/data \
    -disableIPv6"

if [[ "$APP_TOKEN" != "" ]]; then
  CMD_ARGS="$CMD_ARGS -token $APP_TOKEN"
fi

if [[ "$APP_MODE" != "" ]]; then
  CMD_ARGS="$CMD_ARGS -mode $APP_MODE"
fi

if [[ "$APP_METRICS" != "" && "$APP_METRICS" != "false" ]]; then
  CMD_ARGS="$CMD_ARGS -metrics"
fi

if [[ "$APP_URL_SECRET" != "" ]]; then
  CMD_ARGS="$CMD_ARGS -urlSecret $APP_URL_SECRET"
fi

if [[ "$APP_RATE_LIMIT" != "" ]]; then
  CMD_ARGS="$CMD_ARGS -rateLimit $APP_RATE_LIMIT"
fi

if [[ "$APP_CMD_RATE_LIMIT" != "" ]]; then
  CMD_ARGS="$CMD_ARGS -cmdRateLimit $APP_CMD_RATE_LIMIT"
fi

if [[ "$APP_BLACKLIST" != "" ]]; then
  CMD_ARGS="$CMD_ARGS -blacklist $APP_BLACKLIST"
fi

if [[ "$APP_USE_HTTPS" != "" && "$APP_USE_HTTPS" != "false" ]]; then
  CMD_ARGS="$CMD_ARGS -useHttps -certPath $APP_CERT_PATH -keyPath $APP_KEY_PATH"
fi

./telepush $CMD_ARGS $@
