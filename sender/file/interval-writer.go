package file

import (
    "bufio"
    "github.com/qiniu/log"
    "io"
    "sync"
    "time"
)

// 默认的buffer大小
const DEFAULT_BUFFER_SIZE = 32 * 1024 * 1024

// Flush represents a flushed wc.
type Flush struct {
    Path   string
    Writes int64
    Bytes  int64
    Opened time.Time
    Closed time.Time
    Age    time.Duration
}

type intervalWriter struct {
    sync.RWMutex
    buf    *bufio.Writer
    wc     io.WriteCloser
    tick   *time.Ticker
    opened time.Time
    queue  chan *Flush
}

func newIntervalWriter(wc io.WriteCloser, interval time.Duration) (w *intervalWriter, _ error) {

    w = &intervalWriter{}

    if interval != 0 {
        w.tick = time.NewTicker(interval)
        go w.loop()
    }

    w.buf = bufio.NewWriterSize(wc, DEFAULT_BUFFER_SIZE)
    w.opened = time.Now()
    w.wc = wc
    return w, nil
}

func (w *intervalWriter) loop() {
    for range w.tick.C {
        w.Lock()
        w.flush()
        w.Unlock()
    }
}

// Flush for the given reason and re-open.
func (w *intervalWriter) flush() error {
    if w.buf.Size() == 0 {
        return nil
    }

    err := w.buf.Flush()
    if err != nil {
        return err
    }
    return nil
}

func (w *intervalWriter) Close() error {
    w.Lock()
    defer w.Unlock()

    if w.tick != nil {
        w.tick.Stop()
    }
    w.wc.Close()

    return w.flush()
}

func (w *intervalWriter) Write(data []byte) (int, error) {
    w.Lock()
    defer w.Unlock()

    n, err := w.buf.Write(data)
    log.Debugf("[interval writer] write %d bytes, now buffer length: %d bytes, %d bytes available.", len(data), w.buf.Buffered(),w.buf.Available())
    if err != nil {
        return n, err
    }
    return n, err
}
