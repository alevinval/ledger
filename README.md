# Ledger

This library implements an event log, entries can be written, and consumers
of the log will resume the reading since the last known checkpoint. The behaviour
for missing checkpoints can be customized, supporting latest, earliest and custom
offset, each one changing the starting point from which the consumer will read
the log.

It does so with two primitives:

* Ledger writer, writes events and commits offsets.
* Ledger reader, reads written events from a given checkpoint.

The underlying storage is a badger key value store.

Ledger can be used to ensure all events that reach your application
are handled without a loss (in-process), plans for the future could
involve a distributed ledger.

## Development

Running tests with `go test --tags debug ...` will enable debug level logging

## Re-generating protos

Make sure you have `protoc-gen-go` and `protoc-gen-go-grpc`, follow this guide:
https://grpc.io/docs/languages/go/quickstart/

## Project status

This was an internal package for a pet project of mine. Its in process
of being open-sourced and better generalized. Expect breaking changes
in the APIs. It's not a mature project, and the whole idea behind
publishing it is to make it mature.
