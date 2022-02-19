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

type partitionedWriter interface {
	writer
	GetPartitions() int
}

type reader interface {
	Read() (<-chan base.Message, error)
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

func AssertPartitionedWritesAtOffset(t *testing.T, pw partitionedWriter, expectedOffset uint64) {
	for i := 0; i < pw.GetPartitions(); i++ {
		idx, err := pw.Write([]byte("some-data"))
		assert.NoError(t, err)
		assert.Equal(t, expectedOffset, idx)
	}
}

// assert reads with auto-commits between reads.
func AssertReads(t *testing.T, r reader, expected ...string) {
	assertReadsImpl(t, r, true, expected...)
}

// assert reads without auto-commits between reads.
func AssertReadsNoCommit(t *testing.T, r reader, expected ...string) {
	assertReadsImpl(t, r, false, expected...)
}

func assertReadsImpl(t *testing.T, r reader, autoCommit bool, expected ...string) {
	ch, err := r.Read()
	assert.Nil(t, err)
	for i := range expected {
		select {
		case actual, ok := <-ch:
			assert.Equal(t, expected[i], string(actual.Data()))
			if autoCommit && ok {
				actual.Commit()
			}
		case <-time.After(100 * time.Millisecond):
			if expected[i] != "" {
				assert.FailNowf(t, "missing-read", "expected %q, instead read timeout", expected[i])
			}
		}
	}
}
