## Introduction

vol-test is a set of integration tests that is intended to prove and test API support of volume plugins for Docker. vol-test is based upon BATS(https://github.com/sstephenson/bats.git) and depends on some helper libraries - bats-support and bats-assert which are linked as submodules.

vol-test supports testing against remote environments. Remote Docker hosts should have ssh keys configured for access without a password.

## Setup

- Install BATS.

    ```
    git clone https://github.com/sstephenson/bats.git
    cd bats
    sudo ./install.sh /usr/local
    ```

- Clone this repository (optionally, fork), and pull submodules

    ```
    git clone https://github.com/khudgins/vol-test
    cd vol-test
    git submodule init
    git submodule update --recursive --remote
    ```

## Running

- Configuration:

vol-test requires a few environment variables to be configured before running:

* VOLDRIVER - this should be set to the full path (store/vendor/pluginname:tag) of the volume driver to be tested
* PLUGINOPTS - Gets appended to the 'docker volume install' command for install-time plugin configuration
* CREATEOPTS - Optional. Used in 'docker volume create' commands in testing to pass options to the driver being tested
* PREFIX - Optional. Commandline prefix for remote testing. Usually set to 'ssh address_of_node1'
* PREFIX2 - Optional. Commandline prefix for remote testing. Usually set to 'ssh address_of_node2'


- To validate a volume plugin:

1. Export the name of the plugin that is referenced when creating a network as the environmental variable `$VOLDRIVER`.
2. Run the bats tests by running `bats singlenode.bats secondnode.bats`

Example using the vieux/sshfs driver (replace `vieux/sshfs` with the name of the plugin/driver you wish to test):

Prior to running tests the first time, you'll want to pull all the BATS assist submodules, as well:
```
git submodule update --recursive --remote
```

```
$PREFIX="docker-machine ssh node1 "
$VOLDRIVER=vieux/sshfs
$CREATEOPTS="-o khudgins@192.168.99.1:~/tmp -o password=yourpw"

bats singlenode.bats

✓ Test: Create volume using driver (vieux/sshfs)
✓ Test: Confirm volume is created using driver (vieux/sshfs)
...

15 tests, 0 failures
```
