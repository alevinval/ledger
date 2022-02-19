package ledger

type writerListener struct {
	closeManager       chan emptyObj
	closeManagerNotify chan emptyObj
	newReader          chan *Reader
	newWrite           chan emptyObj
	readers            map[string]*Reader
}

func (l *writerListener) notifyWrite() {
	// We don't care if the channel is busy, a burst of writes will be retrieved
	// with a single fetch.
	fireAndForget(l.newWrite)
}

func (l *writerListener) notifyReader(r *Reader) {
	l.newReader <- r
}

func (l *writerListener) close() {
	fireAndWait(l.closeManager, l.closeManagerNotify)
}

func (l *writerListener) Listen() {
	for {
		select {
		case newReader := <-l.newReader:
			// If the reader ID is already tracked, close the tracked reader
			// and replace it for the new instance.
			oldReader, oldReaderPresent := l.readers[newReader.id]
			if oldReaderPresent {
				if oldReader != newReader {
					oldReader.Close()
				}
			}
			l.readers[newReader.id] = newReader
		case <-l.newWrite:
			for _, reader := range l.readers {
				isClosed := reader.doTriggerFetch()
				if isClosed {
					delete(l.readers, reader.id)
				}
			}
		case <-l.closeManager:
			for _, r := range l.readers {
				r.Close()
			}
			l.readers = nil
			close(l.newReader)
			close(l.newWrite)
			close(l.closeManager)
			l.closeManagerNotify <- empty
			close(l.closeManagerNotify)
			return
		}
	}
}
