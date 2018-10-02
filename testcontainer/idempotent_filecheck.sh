#!/bin/sh

if [ -f /data/textfile ]; then
  echo "/data/textfile exists, skipping"
else
  echo "creating /data/textfile with content"
  echo "dockertext" > /data/textfile
fi

if [ -f /data/binaryfile ]; then
  echo "/data/binaryfile exists, skipping"
else
  echo "creating /data/binaryfile"
  dd if=/dev/urandom of=/data/binaryfile bs=10M count=1
  md5sum /data/binaryfile > /data/binchecksum
fi 
