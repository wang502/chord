package chord

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func hashHelper(host string) []byte {
	sh := sha1.New()
	sh.Write([]byte(host))
	bytes := sh.Sum(nil)
	return bytes
}

func TestFindSuccessor(t *testing.T) {
	config1 := DefaultConfig("host1")
	httpTransporter1 := NewTransporter()
	server1 := NewServer("node1", config1, httpTransporter1)
	httpTransporter1.Install(server1, mux.NewRouter())

	config2 := DefaultConfig("host2")
	httpTransporter2 := NewTransporter()
	server2 := NewServer("node2", config2, httpTransporter2)
	httpTransporter2.Install(server2, mux.NewRouter())

	config3 := DefaultConfig("host3")
	httpTransporter3 := NewTransporter()
	server3 := NewServer("node3", config3, httpTransporter3)
	httpTransporter3.Install(server3, mux.NewRouter())

	/* node2 with host "host2" send NewFindSuccessorRequest to node1 with host "host1" */
	// hash("host2") < hash("host1")
	findSuccessorReq1 := NewFindSuccessorRequest(hashHelper("host2"), "host2")
	var data1 bytes.Buffer
	_, err := findSuccessorReq1.Encode(&data1)
	req1, err := http.NewRequest("GET", "/FindSuccessor", &data1)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()

	node1handler := http.HandlerFunc(httpTransporter1.FindSuccessorHandler(server1))
	node1handler.ServeHTTP(rr1, req1)

	findSuccessorResp1 := &FindSuccessorResponse{}
	_, err = findSuccessorResp1.Decode(rr1.Body)
	fmt.Println(findSuccessorResp1.host)
	if err != nil {
		t.Error("response decode error")
	}

	// check response result
	if bytes.Compare([]byte(findSuccessorResp1.ID), server1.node.ID) != 0 || findSuccessorResp1.host != "host1" {
		t.Error("wrong FindSuccessorResponse")
	}

	/* node3 with host "host3" send NewFindSuccessorRequest to node2 with host "host2" */
	// hash("host3") < hash("host2")
	findSuccessorReq2 := NewFindSuccessorRequest(hashHelper("host3"), "host3")
	var data2 bytes.Buffer
	_, err = findSuccessorReq2.Encode(&data2)
	req2, err := http.NewRequest("GET", "/FindSuccessor", &data2)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()

	node2handler := http.HandlerFunc(httpTransporter2.FindSuccessorHandler(server2))
	node2handler.ServeHTTP(rr2, req2)

	findSuccessorResp2 := &FindSuccessorResponse{}
	_, err = findSuccessorResp2.Decode(rr2.Body)
	if err != nil {
		t.Error("response decode error")
	}

	// check response result
	if bytes.Compare([]byte(findSuccessorResp2.ID), server2.node.ID) != 0 || findSuccessorResp2.host != "host2" {
		t.Error("wrong FindSuccessorResponse")
	}

	/* node3 with host "host3" send NewFindSuccessorRequest to node1 with host "host1" */
	// hash("host3") < hash("host1")
	findSuccessorReq3 := NewFindSuccessorRequest(hashHelper("host3"), "host3")
	var data3 bytes.Buffer
	_, err = findSuccessorReq3.Encode(&data3)
	req3, err := http.NewRequest("GET", "/FindSuccessor", &data3)
	if err != nil {
		t.Fatal(err)
	}
	rr3 := httptest.NewRecorder()

	node1handler.ServeHTTP(rr3, req3)

	findSuccessorResp3 := &FindSuccessorResponse{}
	_, err = findSuccessorResp3.Decode(rr3.Body)
	if err != nil {
		t.Error("response decode error")
	}

	// check response result
	if bytes.Compare([]byte(findSuccessorResp3.ID), server1.node.ID) != 0 || findSuccessorResp3.host != "host1" {
		t.Error("wrong FindSuccessorResponse")
	}
}
