package ledger

type writerListener struct {
	readers            map[string]*Reader
	newReader          chan *Reader
	newWrite           chan struct{}
	closeManager       chan struct{}
	closeManagerNotify chan struct{}
}

func (l *writerListener) notifyWrite() {
	// We don't care if the channel is busy, a burst of writes will be retrieved
	// with a single fetch.
	select {
	case l.newWrite <- struct{}{}:
	default:
	}
}

func (l *writerListener) notifyReader(r *Reader) {
	l.newReader <- r
}

func (l *writerListener) close() {
	l.closeManager <- struct{}{}
	<-l.closeManagerNotify
}

func (l *writerListener) manager() {
	for {
		select {
		case r := <-l.newReader:
			// If the reader ID is already tracked, close the tracked reader
			// and replace it for the new instance.
			ret, ok := l.readers[r.id]
			if ok {
				if ret == r {
					continue
				}
				ret.Close()
			}
			l.readers[r.id] = r
		case <-l.newWrite:
			for _, r := range l.readers {
				isClosed := r.doTriggerFetch()
				if isClosed {
					delete(l.readers, r.id)
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
			l.closeManagerNotify <- struct{}{}
			close(l.closeManagerNotify)
			return
		}
	}
}
