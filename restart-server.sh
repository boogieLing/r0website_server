#!/bin/bash
# shellcheck disable=SC1068
pid=$(tail -1 ".pid")
echo $pid
kill -1 $pid

# chmod +x ./restart-server.sh && ./restart-server.sh