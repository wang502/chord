package chord

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestFindSuccessor(t *testing.T) {
	// initialize Chord server, http transporter, router for test
	httpTransporter := NewTransporter()
	server := NewServer("", DefaultConfig("localhost2"), httpTransporter)
	router := mux.NewRouter()
	httpTransporter.Install(server, router)

	// Post find successor request
	findSuccessorReq := NewFindSuccessorRequest([]byte("localhost1"), "localhost1")
	var data bytes.Buffer
	_, err := findSuccessorReq.Encode(&data)
	req, err := http.NewRequest("GET", "/FindSuccessor", &data)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(httpTransporter.FindSuccessorHandler(server))
	handler.ServeHTTP(rr, req)

	findSuccessorResp := &FindSuccessorResponse{}
	_, err = findSuccessorResp.Decode(rr.Body)
	if err != nil {
		t.Error("response decode error")
	}

	// check response result
	if bytes.Compare([]byte(findSuccessorResp.ID), server.node.ID) != 0 || findSuccessorResp.host != "localhost2" {
		t.Error("wrong FindSuccessorResponse")
	}
}
