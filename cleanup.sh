#!/bin/bash

# Cleanup script between test runs.

echo `$PREFIX docker stop mounter`

echo `$PREFIX docker rm mounter`

echo `$PREFIX docker volume rm testvol`

echo `$PREFIX docker plugin disable $VOLDRIVER`

echo `$PREFIX docker plugin rm $VOLDRIVER`

echo `$PREFIX2 docker stop mounter`

echo `$PREFIX2 docker rm mounter`

echo `$PREFIX2 docker volume rm testvol`

echo `$PREFIX2 docker plugin disable $VOLDRIVER`

echo `$PREFIX2 docker plugin rm $VOLDRIVER`
