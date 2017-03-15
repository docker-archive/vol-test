#!/usr/bin/env bats

load test_helper

@test "Test: Install plugin for driver ($driver)" {
  #skip "Faster for rev during development without it - leave driver installed"
  run $prefix docker plugin install --grant-all-permissions $driver $pluginopts
  assert_success
}

@test "Test: Create volume using driver ($driver)" {
  run $prefix docker volume create --driver $driver $createopts testvol
  assert_success
}

@test "Test: Confirm volume is created (volume ls) using driver ($driver)" {
  run $prefix docker volume ls
  assert_line --partial "testvol"
}

@test "Test: Confirm docker volume inspect works using driver ($driver)" {
  run $prefix docker volume inspect testvol
  assert_line --partial "\"Driver\": \"$driver"
}

@test "Start a container and mount the volume" {
  run $prefix docker run -it -d --name mounter -v testvol:/data ubuntu /bin/bash
  assert_success
}

@test "Write a textfile to the volume" {
  run $prefix -t 'docker exec -it mounter /bin/bash -c "echo \"testdata\" > /data/foo.txt"'
  assert_success
}

@test "Confirm textfile contents on the volume" {
  run $prefix -t docker exec -it mounter cat /data/foo.txt
  assert_line --partial "testdata"
}

@test "Create a binary file" {
  run $prefix -t docker exec -it 'mounter dd if=/dev/urandom of=/data/random bs=10M count=1'
  assert_output --partial "10 M"
}

@test "get a checksum for that binary file" {
  run $prefix -t 'docker exec -it mounter /bin/bash -c "md5sum /data/random > /data/checksum"'
  assert_success
}

@test "Confirm checksum" {
  run $prefix -t docker exec -it mounter md5sum --check /data/checksum
  assert_success
}

@test "Destroy container" {
  run $prefix docker stop mounter
  run $prefix docker rm mounter
  assert_success
}

@test "Let's see if our data is still here" {
  run $prefix docker run -it -d --name mounter -v testvol:/data ubuntu /bin/bash
  assert_success
}

@test "Confirm textfile contents are still on the volume" {
  run $prefix -t docker exec -it mounter cat /data/foo.txt
  assert_line --partial "testdata"
}

@test "Confirm checksum persistence" {
  run $prefix -t docker exec -it mounter md5sum --check /data/checksum
  assert_success
}

@test "Destroy container" {
  run $prefix docker stop mounter
  run $prefix docker rm mounter
  assert_success
}
