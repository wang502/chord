package chord

import (
	"bytes"
	"net/http"

	"log"

	"fmt"

	"time"

	"github.com/gorilla/mux"
)

// Transporter represents a http communication gate with other nodes
type Transporter struct {
	httpClient         http.Client
	listNodesPath      string
	findSuccessorPath  string
	notifyPath         string
	getPredecessorPath string
	getSuccessorPath   string
	setPredecessorPath string
	joinPath           string
	startPath          string
	stopPath           string

	getFingerTablePath string
}

// NewTransporter initilizes a new Transporter object
func NewTransporter() *Transporter {
	return &Transporter{
		httpClient:         http.Client{Timeout: time.Second},
		findSuccessorPath:  "/findSuccessor",
		notifyPath:         "/notify",
		getPredecessorPath: "/getPredecessor",
		getSuccessorPath:   "/getSuccessor",
		joinPath:           "/join",
		startPath:          "/start",
		stopPath:           "/stop",

		getFingerTablePath: "/getFingerTable",
	}
}

// Install applies the chord route to an http router
func (t *Transporter) Install(server *Server, mux *mux.Router) {
	mux.HandleFunc(t.notifyPath, t.notifyHandler(server))
	mux.HandleFunc(t.findSuccessorPath, t.findSuccessorHandler(server))
	mux.HandleFunc(t.getPredecessorPath, t.getPredecessorHandler(server))
	mux.HandleFunc(t.getSuccessorPath, t.getSuccessorHandler(server))
	mux.HandleFunc(t.joinPath, t.joinHandler(server)).Methods("POST")
	mux.HandleFunc(t.startPath, t.startHandler(server)).Methods("POST")
	mux.HandleFunc(t.stopPath, t.stopHandler(server)).Methods("POST")

	mux.HandleFunc(t.getFingerTablePath, t.getFingerTableHandler(server))
}

// -------------------------------------------------------------------------
//
// Sending request
//
// -------------------------------------------------------------------------

// SendFindSuccessorRequest sends outgoing find successor request to other Node server, a successor response will be returned
func (t *Transporter) SendFindSuccessorRequest(server *Server, req *FindSuccessorRequest) (*FindSuccessorResponse, error) {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		return nil, err
	}

	url := req.host + t.findSuccessorPath
	resp, err := t.httpClient.Post(url, "chord.protobuf", &b)

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

//	-------------------------------------------------------------------------
//
//	handler functions
//
//	-------------------------------------------------------------------------

// findSuccessorHandler handles incoming request to find successor of the given key
func (t *Transporter) findSuccessorHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &FindSuccessorRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp, err := server.FindSuccessor(req)
		if resp == nil || err != nil {
			http.Error(w, "Failed to return successor information", http.StatusBadRequest)
			return
		}

		log.Printf("host %s's successor is: %s", req.host, resp.host)
		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// notifyHandler handles incoming notify about possibe new predecessor
func (t *Transporter) notifyHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &NotifyRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp, err := server.notify(req)
		if resp == nil || err != nil {
			http.Error(w, "failed to notify", http.StatusBadRequest)
			return
		}

		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// getPredecessorHandler handles incoming request to return this local server's predecessor
func (t *Transporter) getPredecessorHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		predResp, err := server.handleGetPredecessorRequest()
		if predResp == nil || err != nil {
			http.Error(w, "failed to return predecessor", http.StatusBadRequest)
			return
		}
		log.Printf("host %s's predecessor is %s", server.config.Host, predResp.host)

		if _, err := predResp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// getSuccessorHandler handles the incoming request to return this node's successor
func (t *Transporter) getSuccessorHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		succResp, err := server.handleGetSuccessorRequest()
		if err != nil {
			http.Error(w, "failed to return successor", http.StatusBadRequest)
			return
		}

		log.Printf("host %s's successor is %s", server.config.Host, succResp.host)

		if _, err := succResp.Encode(w); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

// joinHandler handles the post request for this server to join an existing Chord ring
// the url pattern is '/join?host='
func (t *Transporter) joinHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.URL.Query().Get("host")
		log.Println(host)
		err := server.Join(host)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to join %s.%s", host, err), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "success to join %s", host)
	}
}

// startHandler handles the incoming request to start this Chord server
func (t *Transporter) startHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := server.Start()
		if err != nil {
			fmt.Fprintf(w, "error to start server %s", server.config.Host)
		} else {
			fmt.Fprintf(w, "success to start server %s", server.config.Host)
		}
	}
}

// stopHandler handles incoming request to stop the Chord server
func (t *Transporter) stopHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := server.Stop()
		if err != nil {
			fmt.Fprintf(w, "error to stop server %s", server.config.Host)
		} else {
			fmt.Fprintf(w, "success to stop server %s", server.config.Host)
		}
	}
}

// getFingerTableHandler handles incoming request to log entries in finger table
func (t *Transporter) getFingerTableHandler(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fingers := server.node.Finger()
		for i := 0; i < server.config.HashBits; i++ {
			log.Printf("host %s's finger at index %d: %s", server.config.Host, i, fingers[i])
		}
	}
}
