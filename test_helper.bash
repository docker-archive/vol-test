#!/usr/bin/env bats

load 'test/test_helper/bats-support/load'
load 'test/test_helper/bats-assert/load'

driver=$VOLDRIVER
prefix=$PREFIX
createopts="$CREATEOPTS"
