#!/usr/bin/env bats

load 'test/test_helper/bats-support/load'
load 'test/test_helper/bats-assert/load'

driver=$VOLDRIVER
prefix=$PREFIX
prefix2=$PREFIX2
createopts="$CREATEOPTS"
pluginopts="$PLUGINOPTS"
