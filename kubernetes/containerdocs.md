## Container Details

The test container is pretty simple: it's a Ruby-based webapp built with Sinatra that implements an extremely simple, lightweight HTTP api. It's not REST - everything is done via simple GET requsts.

The container's spec is in the /testcontainer directory in this repository. You can build it yourself, or pull it from Docker Hub at khudgins/volcheck

## Container API

The container API has a few methods, all available via HTTP GET statements. This is NOT a rest API, just something simple for testing.

# /resetfilecheck

Resets the test data to begin a clean test run

# /runfilecheck

Creates the datafiles needed to perform volume function tests. We write a known phrase into a textfile, and randomly write a binary file. Then we run an md5sum against the binary file and store the results on disk in the test volume.

# /textcheck

Returns "1" if the test textfile contains the correct, known data. Returns "0" if not, or if the file doesn't exist.

# /bincheck

Returns "1" if the test binaryfile matches its original checksum. Returns "0" if not (including if the file does not exist)

# /status

Returns "OK" as a container healthcheck. This is hard-coded - if the container is running, you'll get an "OK".

# /shutdown

Immediately terminates the container process. This will trigger Kubernetes to spawn a replacement container. The container will respawn on the same node as its previous incarnation unless you cordon the node or some other method of forcing a migration, like kubectl drain.
