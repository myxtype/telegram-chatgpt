#!/bin/bash

cmds=$(ls -l ./ | awk '/^d/ {print $NF}')

for d in $cmds; do
  cd $d || exit
  echo "building ${d} ..."
  go build -tags=jsoniter .
  cd ../ || exit
done
