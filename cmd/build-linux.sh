#!/bin/bash

export GOOS=linux

cmds=$(ls -l ./ | awk '/^d/ {print $NF}')

for d in $cmds; do
  cd $d || exit
  echo "building ${d} ..."
  go build -tags=jsoniter .
  cd ../ || exit
done
