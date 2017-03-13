package chord

import (
	"bytes"
	"net/http"

	"log"

	"github.com/gorilla/mux"
)

// Transporter represents a http communication gate with other nodes
type Transporter struct {
	httpClient         http.Client
	listNodesPath      string
	findSuccessorPath  string
	notifyPath         string
	getPredecessorPath string
	setPredecessorPath string
}

// NewTransporter initilizes a new Transporter object
func NewTransporter() *Transporter {
	return &Transporter{
		httpClient:         http.Client{},
		findSuccessorPath:  "/findSuccessor",
		notifyPath:         "/notify",
		getPredecessorPath: "/getPredecessor",
	}
}

// Install applies the chord route to an http router
func (t *Transporter) Install(server *Server, mux *mux.Router) {
	mux.HandleFunc(t.notifyPath, t.NotifyHandler(server))
	mux.HandleFunc(t.findSuccessorPath, t.FindSuccessorHandler(server))
	mux.HandleFunc(t.getPredecessorPath, t.GetPredecessorHandler(server))
}

// Sending

// SendFindSuccessorRequest sends outgoing find successor request to other Node server, a successor response will be returned
func (t *Transporter) SendFindSuccessorRequest(server *Server, req *FindSuccessorRequest) (*FindSuccessorResponse, error) {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		return nil, err
	}

	url := req.host + t.findSuccessorPath
	resp, err := t.httpClient.Post(url, "chord.protobuf", &b)

	//data, _ := ioutil.ReadAll(resp.Body)
	//log.Println(len(data))

	if err != nil {
		return nil, err
	}

	successorResp := &FindSuccessorResponse{}
	if _, err = successorResp.Decode(resp.Body); err != nil {
		return nil, err
	}

	return successorResp, nil
}

// SendNotifyRequest sends a request to other node to nofify it about the possible new predecessor
func (t *Transporter) SendNotifyRequest(server *Server, req *NotifyRequest) (*NotifyResponse, error) {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		return nil, err
	}

	url := req.targetHost + t.notifyPath
	resp, err := t.httpClient.Post(url, "chord.protobuf", &b)

	if err != nil {
		return nil, err
	}

	notifyResp := &NotifyResponse{}
	if _, err = notifyResp.Decode(resp.Body); err != nil {
		return nil, err
	}

	return notifyResp, nil
}

// SendGetPredecessorRequest sends a request to get the predecessor of server on given host
func (t *Transporter) SendGetPredecessorRequest(server *Server, host string) (*GetPredecessorResponse, error) {
	url := host + t.getPredecessorPath
	resp, err := t.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	predResp := &GetPredecessorResponse{}
	if _, err = predResp.Decode(resp.Body); err != nil {
		return nil, err
	}

	return predResp, nil
}

// Receiving

// FindSuccessorHandler handles incoming request to find successor of the given key
func (t *Transporter) FindSuccessorHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &FindSuccessorRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp, err := server.FindSuccessor(req)
		log.Printf("succ: %s", resp.host)
		if resp == nil || err != nil {
			http.Error(w, "Failed to return successor information", http.StatusBadRequest)
			return
		}

		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// NotifyHandler handles incoming notify about possibe new predecessor
func (t *Transporter) NotifyHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &NotifyRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp, err := server.Notify(req)
		if resp == nil || err != nil {
			http.Error(w, "failed to return successor", http.StatusBadRequest)
			return
		}

		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// GetPredecessorHandler handles incoming request to return this local server's predecessor
func (t *Transporter) GetPredecessorHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		predResp, err := server.HandleGetPredecessorRequest()
		if predResp == nil || err != nil {
			http.Error(w, "failed to return predecessor", http.StatusBadRequest)
			return
		}
		log.Printf("pred: %s", predResp)

		if _, err := predResp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}
