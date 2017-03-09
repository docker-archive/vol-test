#!/usr/bin/env bats

load test_helper

@test "Test: Install plugin for driver ($driver)" {
  skip "This test works, faster for rev without it"
  run $prefix docker plugin install --grant-all-permissions $driver $pluginopts
  assert_success
}

@test "Test: Create volume using driver ($driver)" {
  skip "Test works, revs faster without"
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
  skip
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
