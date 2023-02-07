#!/bin/sh
set -x
for port in $PORTS; do
    /usr/local/bin/port-tester --proto $PROTO --port $port --sleep $SLEEP --delay 2 $TARGETS &
done
wait
