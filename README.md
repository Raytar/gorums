# gorums

Gorums [1] is a framework for simplifying the design and implementation of
fault-tolerant quorum-based protocols. Gorums allows to group replicas into one
or more _configurations_. A configuration also holds information on how many
replicas are necessary to form a quorum. Gorums enables programmers to invoke
remote procedure calls (RPCs) on the replicas in a configuration and wait for
responses from a quorum. We call this a quorum remote procedure call (QRPC).

Gorums uses code generation to produce an RPC library that clients can use to
invoke QRPCs. Gorums is a wrapper around the [gRPC](http://www.grpc.io/)
library. Services are defined using the protocol buffers interface definition
language.

### Examples

A collection of different algorithms for reconfigurable atomic storage
implemented using Gorums can be found
[here](https://github.com/relab/smartmerge).

### Benchmarking

TODO Clean this up:

Ca. slik for en spesifik bencmark:
Trenger: go get -u  golang.org/x/tools/cmd/benchcmp

cd dev
git checkout gorumsgen

For nyeste commit med endringer:
go test . -run=NONE -bench=Read1KQ1N1Local -benchtime=5s > new.txt

Så bytte til forrige commit:
git checkout HEAD~1

go test . -run=NONE -bench=Read1KQ1N1Local -benchtime=5s > old.txt

Compare:
benchcmp old.txt new.txt

### Documentation

* [Student/user guide](doc/userguide.md)
* [Developer guide](doc/devguide.md)

### References

[1] Tormod E. Lea, Leander Jehl, and Hein Meling. Gorums: _A Framework for
    Implementing Reconfigurable Quorum-based Systems._ In submission.
