#!/bin/sh
echo "$@" | grep server
if [ $? -ne 0 ]; then
  echo "start cloudprober service"
  cloudprober -config_file cloudprober.cfg &
  sleep 2
fi
echo "start service prober"
echo "executing: /app/service-prober $@"

/app/service-prober $@