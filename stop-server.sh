#!/bin/bash
# shellcheck disable=SC1068
pid=$(tail -1 ".pid")
echo $pid
kill -9 $pid

# chmod +x ./stop-server.sh && ./stop-server.sh