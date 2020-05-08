package file

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

type wcTest struct {
    *bytes.Buffer
}

func (w wcTest) Close() error {
    return nil
}
func Test_intervalWriter(t *testing.T) {
    result := wcTest{
        Buffer: bytes.NewBuffer([]byte{}),
    }
    iw, err := newIntervalWriter(result, 2*time.Second)

    assert.Nil(t, err)
    // 开始并发测试
    //var wg sync.WaitGroup

    _, err = iw.Write([]byte("this is a test"))
    assert.Nil(t, err)

    assert.True(t, result.Len() == 0)
    time.Sleep(3 * time.Second)
    assert.True(t, result.Len() > 0)

}
