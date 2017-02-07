## Setup

- Install BATS.

    ```
    git clone https://github.com/sstephenson/bats.git
    cd bats
    sudo ./install.sh /usr/local
    ```

## Running

- To validate a volume plugin:

1. Export the name of the plugin that is referenced when creating a network as the environmental variable `$VOLDRIVER`.
2. Run the bats tests.

Example using the vieux/sshfs driver (replace `vieux/sshfs` with the name of the plugin/driver you wish to test):

```
$PREFIX="docker-machine ssh node1 "
$VOLDRIVER=vieux/sshfs
$CREATEOPTS="-o khudgins@192.168.99.1:~/tmp -o password=yourpw"

bats build.bats

✓ Test: Create volume using driver (vieux/sshfs)
✓ Test: Confirm volume is created using driver (vieux/sshfs)

2 tests, 0 failures
```
