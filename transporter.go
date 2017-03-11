package chord

import (
	"bytes"
	"net/http"

	"log"

	"github.com/gorilla/mux"
)

// Transporter represents a http communication gate with other nodes
type Transporter struct {
	httpClient        http.Client
	listNodesPath     string
	findSuccessorPath string
	notifyPath        string
}

// NewTransporter initilizes a new Transporter object
func NewTransporter() *Transporter {
	return &Transporter{
		httpClient:        http.Client{},
		listNodesPath:     "/listNodes",
		findSuccessorPath: "/findSuccessor",
		notifyPath:        "/notify",
	}
}

// Install applies the chord route to an http router
func (t *Transporter) Install(server *Server, mux *mux.Router) {
	mux.HandleFunc(t.listNodesPath, t.ListNodesHandler(server))
	mux.HandleFunc(t.findSuccessorPath, t.FindSuccessorHandler(server))
}

// Sending

// SendListNodesRequest sends a list nodes request to other Node server, a list nodes response will be returned
func (t *Transporter) SendListNodesRequest(server *Server, req *ListNodesRequest, existing string) *ListNodesResponse {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		log.Println("chord.ListNodes.encoding.error")
		return nil
	}

	url := existing + t.listNodesPath
	resp, err := t.httpClient.Post(url, "chord.protobuf", &b)
	if err != nil {
		log.Println("chord.ListNodes.response.error")
		return nil
	}

	listNodesResp := &ListNodesResponse{}
	if _, err = listNodesResp.Decode(resp.Body); err != nil {
		log.Println("chord.ListNodes.decoding.error")
		return nil
	}

	return listNodesResp
}

// SendFindSuccessorRequest sends outgoing find successor request to other Node server, a successor response will be returned
func (t *Transporter) SendFindSuccessorRequest(server *Server, req *FindSuccessorRequest) *FindSuccessorResponse {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		log.Println("chord.FindSuccessor.encoding.error")
		return nil
	}

	url := req.host + t.findSuccessorPath
	resp, err := t.httpClient.Post(url, "chord.protobuf", &b)
	if err != nil {
		log.Println("chord.FindSuccessor.response.error")
		return nil
	}

	successorResp := &FindSuccessorResponse{}
	if _, err = successorResp.Decode(resp.Body); err != nil {
		log.Println("chord.FindSuccessor.decoding.error")
		return nil
	}

	return successorResp
}

// SendNotifyRequest sends a request to other node to nofify it about the possible new predecessor
func (t *Transporter) SendNotifyRequest(server *Server, req *NotifyRequest) *NotifyResponse {
	return nil
}

// Receiving

// ListNodesHandler handles incoming request to list all nodes
func (t *Transporter) ListNodesHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &ListNodesRequest{}

		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := server.ListNodes(req)
		if resp == nil {
			http.Error(w, "Failed to return nodes information", http.StatusBadRequest)
			return
		}

		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// FindSuccessorHandler handles incoming request to find successor of the given key
func (t *Transporter) FindSuccessorHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &FindSuccessorRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := server.FindSuccessor(req)
		if resp == nil {
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

	}
}
