#!/bin/sh
echo "$@" | grep server
if [ $? -ne 0 ]; then
  echo "start cloudprober service"
  cloudprober -config_file cloudprober.cfg &
fi

# wait a few seconds before start
echo wait 3 seconds...
sleep 3

echo "start service prober"
echo "executing: /app/service-prober $@"

/app/service-prober $@