package chord

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestTransporter(t *testing.T) {
	server := NewServer("", DefaultConfig("localhost"))
	httpTransporter := NewTransporter()
	router := mux.NewRouter()
	httpTransporter.Install(server, router)

	findSuccessorReq := NewFindSuccessorRequest([]byte("test"), "localhost")
	var data bytes.Buffer
	_, err := findSuccessorReq.Encode(&data)
	req, err := http.NewRequest("GET", "/FindSuccessor", &data)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(httpTransporter.FindSuccessorHandler(server))
	handler.ServeHTTP(rr, req)
}
