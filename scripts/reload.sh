#!/bin/bash

# This is a janky bash script to restart the server process when we
# edit files. You should run it from the repo root.
go run *.go &
while inotifywait -e modify *.go; do
    ps aux | grep 'go.*client' | awk '{print $2}' | xargs kill -SIGINT
    go run *.go &
done
