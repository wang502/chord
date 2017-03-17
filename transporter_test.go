package chord

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestFindSuccessorHandler(t *testing.T) {
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

	handler := http.HandlerFunc(httpTransporter.findSuccessorHandler(server))
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

// To run this test, go to /example folder and start a HTTP server in the backgraound at port 1000 and port 2000
// started server acts as an exsiting host in Chord ring
func TestJoinHalder(t *testing.T) {
	resp, err := http.Post("http://localhost:2000/join?host=http://localhost:2001", "chord.join", nil)
	buf := resp.Body
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "success to join http://localhost:2001" {
		log.Println(string(data))
		t.Errorf("failted to join")
	}
}
