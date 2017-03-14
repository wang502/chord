package chord

import (
	"bytes"
	"testing"
)

func createNotifyRequest() (*NotifyRequest, []byte) {
	notifyReq := NewNotifyRequest([]byte("http://localhost:3000"), "http://localhost:3000", "http://localhost:4000")

	var buf bytes.Buffer
	notifyReq.Encode(&buf)
	return notifyReq, buf.Bytes()
}

func BenchmarkNotifyRequestEncoding(b *testing.B) {
	notifReq, _ := createNotifyRequest()
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		notifReq.Encode(&buf)
	}
}

func BenchmarkNotifyRequestDecoding(b *testing.B) {
	notifReq, data := createNotifyRequest()
	for i := 0; i < b.N; i++ {
		notifReq.Decode(bytes.NewReader(data))
	}
}
