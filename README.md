# Ledger

This library implements an event log and imitates how kafka consumers work.
Writers commit entries into the log, and readers consume events. Each reader
can be configured to consume from the first offset, from the last offset
or from a custom offset. Once a reader has a committed offset, when it is
stopped/restarted, the reader will resume from the last committed offset.

These are the basic building blocks:

* `ledger.Writer` writes events to the log
* `ledger.Reader` reads events from the event log
* `ledger.Message` returned when reading, when the message is
processed trigger a commit with `msg.Commit()` to advance the reader offset

Additionally, there are:

* `ledger.PartitionedWriter` which internally keeps as many logs
as partitions. Uses a round-robin strategy to distribute the writes.
* `ledger.PartitionedReader` which supports reading messages from
the partitions, emitting them in the original write order

The underlying storage is a badger key value store. Badger can be
configured to persist on disk, or keep everything in memory.

## Development

Run tests with `go test --tags debug ...` to enable debug level logging

## Re-generating protos

Make sure you have `protoc-gen-go` and `protoc-gen-go-grpc`, follow this guide:
https://grpc.io/docs/languages/go/quickstart/

Then run `make compile-protos`
