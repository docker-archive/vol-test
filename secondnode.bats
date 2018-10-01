#!/usr/bin/env bats

load test_helper

@test "Test: Install plugin for driver ($driver) on node 2" {
  #skip "This test works, faster for rev without it"
  run $prefix2 docker plugin install --grant-all-permissions $driver $pluginopts
  assert_success
}

@test "Test: Confirm volume is visible on second node (volume ls) using driver ($driver)" {
  run $prefix2 docker volume ls
  assert_line --partial "testvol"
}

@test "Start a container and mount the volume on node 2" {
  run $prefix2 docker run -it -d --name mounter -v testvol:/data ubuntu /bin/bash
  assert_success
}

@test "Confirm textfile contents on the volume from node 2" {
  run $prefix2 -t docker exec -it mounter cat /data/foo.txt
  assert_line --partial "testdata"
}

@test "Confirm checksum for binary file on node 2" {
  run $prefix2 -t docker exec -it mounter md5sum --check /data/checksum
  assert_success
}

@test "Destroy container on node 2" {
  run $prefix2 docker stop mounter
  run $prefix2 docker rm mounter
  assert_success
}

@test "Remove volume" {
  run $prefix2 docker volume rm testvol
  assert_success
}

@test "Confirm volume is removed from docker ls" {
  run $prefix2 docker volume ls
  refute_output --partial 'testvol'
}

@test "Disable plugin on node 2" {
  run $prefix2 docker plugin disable $driver
  assert_success
}

@test "Remove plugin on node 2" {
  run $prefix2 docker plugin rm $driver
  assert_success
}

@test "Disable plugin on node 1" {
  run $prefix1 docker plugin disable $driver
  assert_success
}

@test "Remove plugin on node 1" {
  run $prefix1 docker plugin rm $driver
  assert_success
}
