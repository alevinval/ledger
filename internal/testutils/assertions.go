package testutils

import (
	"testing"
	"time"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/stretchr/testify/assert"
)

type checkpointable interface {
	GetCheckpoint() (*proto.Checkpoint, error)
}

type writer interface {
	Write([]byte) (uint64, error)
}

type reader interface {
	Read() (<-chan base.Message, error)
	Commit(uint64) error
}

type partitionedReader interface {
	Read() (<-chan base.PartitionedMessage, error)
}

func AssertCheckpointAt(t *testing.T, c checkpointable, expected uint64) {
	cp, err := c.GetCheckpoint()

	assert.NoError(t, err)
	assert.Equal(t, expected, cp.Offset)
}

func AssertWrites(t *testing.T, w writer, data ...string) {
	for _, payload := range data {
		t.Logf("writing %q\n", payload)
		_, err := w.Write([]byte(payload))
		assert.NoError(t, err)
	}
	time.Sleep(100 * time.Millisecond)
}

// assert reads with auto-commits between reads.
func AssertReads(t *testing.T, r reader, expected ...string) {
	assertReadsImpl(t, r, true, expected...)
}

// assert reads without auto-commits between reads.
func AssertReadsNoCommit(t *testing.T, r reader, expected ...string) {
	assertReadsImpl(t, r, false, expected...)
}

// assert reads with auto-commits between reads.
func AssertReadsPartitioned(t *testing.T, r partitionedReader, expected ...string) {
	assertReadsPartitionedImpl(t, r, true, expected...)
}

func assertReadsImpl(t *testing.T, r reader, autoCommit bool, expected ...string) {
	ch, err := r.Read()
	assert.Nil(t, err)
	for i := range expected {
		select {
		case actual := <-ch:
			assert.Equal(t, expected[i], string(actual.Data()))
			if autoCommit {
				r.Commit(actual.Offset())
			}
		case <-time.After(100 * time.Millisecond):
			if expected[i] != "" {
				assert.FailNowf(t, "missing-read", "expected %q, instead read timeout", expected[i])
			}
		}
	}
}

// TODO: when generics come out, we can merge assertReadsImpl and assertReadPartitionedImpl
func assertReadsPartitionedImpl(t *testing.T, r partitionedReader, autoCommit bool, expected ...string) {
	ch, err := r.Read()
	assert.Nil(t, err)
	for i := range expected {
		select {
		case actual := <-ch:
			assert.Equal(t, expected[i], string(actual.Data()))
			if autoCommit {
				err := actual.Commit()
				assert.NoError(t, err)
			}
		case <-time.After(100 * time.Millisecond):
			if expected[i] != "" {
				assert.FailNowf(t, "missing-read", "expected %q, instead read timeout", expected[i])
			}
		}
	}
}
