#!/usr/bin/env bats

load test_helper

@test "Test: Create volume using driver ($driver)" {
  echo $prefix
  echo $createopts
  run $prefix docker volume create --driver $driver $createopts testvol
  [ "$status" -eq 0 ]
}

@test "Test: Confirm volume is created using driver ($driver)" {
  run $prefix docker volume ls
  assert_line "$driver:latest   testvol"

}
