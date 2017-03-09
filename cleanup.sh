#!/bin/bash

# Cleanup script between test runs.

echo `$PREFIX docker stop mounter`

echo `$PREFIX docker rm mounter`

echo `$PREFIX docker rm testvol`

#echo `$PREFIX docker plugin disable $VOLDRIVER`

#echo `$PREFIX docker plugin rm $VOLDRIVER`
