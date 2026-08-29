#!/bin/sh
set -eu

pids=""

cleanup() {
  trap - EXIT INT TERM
  for pid in $pids; do
    kill "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
}

trap cleanup EXIT INT TERM

make run-administration &
pids="$pids $!"
make run-user &
pids="$pids $!"
make run-ui &
pids="$pids $!"

wait

