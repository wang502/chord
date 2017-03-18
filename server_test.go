package chord

import (
	"bytes"
	"math/big"
	"net/http"
	"testing"
	"time"

	"log"
)

func hashHelper(host string) []byte {
	config := DefaultConfig("")
	hash := config.HashFunc
	hash.Write([]byte(host))
	b := hash.Sum(nil)

	idInt := big.Int{}
	idInt.SetBytes(b)

	// Get the ceiling
	two := big.NewInt(2)
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(config.HashBits)), nil)

	// Apply the mod
	idInt.Mod(&idInt, &ceil)

	// Add together
	return idInt.Bytes()
}

func TestStartI(t *testing.T) {
	// initialize an http client for sending request to test servers
	client := http.Client{}

	// Test server1 joins an exsiting Chord ring, from exisiting host "http://localhost:3000"
	_, err := client.Post("http://localhost:5000/join?host=http://localhost:6000", "chord.join", nil)

	// Test server1's start function, check whether it periodically stabilizes
	_, err = client.Post("http://localhost:5000/start", "chord.start", nil)
	_, err = client.Post("http://localhost:6000/start", "chord.start", nil)

	time.Sleep(5 * DefaultStabilizeInterval)

	resp1, _ := client.Get("http://localhost:6000/getPredecessor")
	predResp1 := &GetPredecessorResponse{}
	_, err = predResp1.Decode(resp1.Body)
	if err != nil {
		t.Errorf("error decode.%s", err)
	}
	if bytes.Compare([]byte(predResp1.ID), hashHelper("http://localhost:5000")) != 0 || predResp1.host != "http://localhost:5000" {
		t.Errorf("wrong predecessor returned")
	}

	resp2, _ := client.Get("http://localhost:6000/getSuccessor")
	succResp2 := &FindSuccessorResponse{}
	_, err = succResp2.Decode(resp2.Body)
	if err != nil {
		t.Errorf("error decode.%s", err)
	}
	if bytes.Compare([]byte(succResp2.ID), hashHelper("http://localhost:5000")) != 0 || succResp2.host != "http://localhost:5000" {
		t.Errorf("wrong successor returned")
	}

	resp3, _ := client.Get("http://localhost:5000/getPredecessor")
	predResp3 := &GetPredecessorResponse{}
	_, err = predResp3.Decode(resp3.Body)
	if err != nil {
		t.Errorf("error decode.%s", err)
	}
	if predResp3.host != "http://localhost:6000" {
		t.Errorf("wrong predecessor returned")
	}

	resp4, _ := client.Get("http://localhost:5000/getSuccessor")
	succResp4 := &FindSuccessorResponse{}
	_, err = succResp4.Decode(resp4.Body)
	if err != nil {
		t.Errorf("error decode.%s", err)
	}
	if succResp4.host != "http://localhost:6000" {
		log.Printf("[TEST] %s", succResp4.host)
		t.Errorf("wrong successor returned")
	}

	// -------------------------
	//
	// -------------------------
	// Test server1 joins an exsiting Chord ring, from exisiting host "http://localhost:3000"
	//
	_, err = client.Post("http://localhost:7000/join?host=http://localhost:5000", "chord.join", nil)
	_, err = client.Post("http://localhost:7000/start", "chord.start", nil)

	time.Sleep(5 * DefaultStabilizeInterval)

	resp5, _ := client.Get("http://localhost:7000/getPredecessor")
	predResp5 := &GetPredecessorResponse{}
	_, err = predResp5.Decode(resp5.Body)
	if err != nil {
		t.Errorf("error decode.%s", err)
	}
	if bytes.Compare([]byte(predResp5.ID), hashHelper("http://localhost:5000")) != 0 || predResp5.host != "http://localhost:5000" {
		t.Errorf("wrong predecessor returned")
	}

	resp6, _ := client.Get("http://localhost:7000/getSuccessor")
	predResp6 := &GetPredecessorResponse{}
	_, err = predResp6.Decode(resp6.Body)
	if err != nil {
		t.Errorf("error decode.%s", err)
	}
	if predResp6.host != "http://localhost:6000" {
		t.Errorf("wrong predecessor returned")
	}

	time.Sleep(5 * DefaultStabilizeInterval)

	_, err = client.Get("http://localhost:5000/getFingerTable")
	_, err = client.Get("http://localhost:6000/getFingerTable")
	_, err = client.Get("http://localhost:7000/getFingerTable")

	_, err = client.Post("http://localhost:5000/stop", "chord.stop", nil)
	_, err = client.Post("http://localhost:6000/stop", "chord.stop", nil)
	_, err = client.Post("http://localhost:7000/stop", "chord.stop", nil)
}

// Test for 3 servers, the hashed bytes are in an order as belowed
// 8000 < 10000 < 9000
func TestStartII(t *testing.T) {
	client := http.Client{}

	_, err := client.Post("http://localhost:9000/join?host=http://localhost:8000", "chord.join", nil)
	if err != nil {
		t.Errorf("failed to join, %s", err)
	}
	_, err = client.Post("http://localhost:9000/start", "chord.start", nil)
	_, err = client.Post("http://localhost:8000/start", "chord.start", nil)

	_, err = client.Post("http://localhost:10000/join?host=http://localhost:8000", "chord.join", nil)
	_, err = client.Post("http://localhost:10000/start", "chord.start", nil)

	time.Sleep(5 * DefaultStabilizeInterval)

	/*_, err = client.Get("http://localhost:8000/getFingerTable")
	_, err = client.Get("http://localhost:10000/getFingerTable")
	_, err = client.Get("http://localhost:9000/getFingerTable")*/

	_, err = client.Post("http://localhost:9500/join?host=http://localhost:9000", "chord.join", nil)
	if err != nil {
		t.Errorf("failed to join, %s", err)
	}
	_, err = client.Post("http://localhost:9500/start", "chord.start", nil)

	time.Sleep(15 * DefaultStabilizeInterval)

	_, err = client.Get("http://localhost:9500/getFingerTable")
	_, err = client.Get("http://localhost:8000/getFingerTable")
	_, err = client.Get("http://localhost:10000/getFingerTable")
	_, err = client.Get("http://localhost:9000/getFingerTable")

	_, err = client.Post("http://localhost:8000/stop", "chord.stop", nil)
	_, err = client.Post("http://localhost:9000/stop", "chord.stop", nil)
	_, err = client.Post("http://localhost:10000/stop", "chord.stop", nil)
	_, err = client.Post("http://localhost:9500/stop", "chord.stop", nil)
}
