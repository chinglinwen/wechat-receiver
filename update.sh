#!/bin/bash

# set -e

. ~/bin/k8s

name=wechat-rec

sh build-image.sh
8del $name
while true; do
  out="$( 8pods $name )"
  echo "$out"
  echo "$out" | grep Run >/dev/null && break
  sleep 1
done