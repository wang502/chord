package chord

import (
	"bytes"
	"net/http"

	"log"

	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Transporter represents a http communication gate with other nodes
type Transporter struct {
	httpClient        http.Client
	listNodesPath     string
	findSuccessorPath string

	getPredecessorPath string
	getSuccessorPath   string
	setPredecessorPath string
	getFingerTablePath string

	notifyPath string
	joinPath   string
	startPath  string
	stopPath   string
}

// NewTransporter initilizes a new Transporter object
func NewTransporter() *Transporter {
	return &Transporter{
		httpClient:         http.Client{Timeout: time.Second},
		findSuccessorPath:  "/findSuccessor",
		getPredecessorPath: "/getPredecessor",
		getSuccessorPath:   "/getSuccessor",
		getFingerTablePath: "/getFingerTable",
		notifyPath:         "/notify",
		joinPath:           "/join",
		startPath:          "/start",
		stopPath:           "/stop",
	}
}

// Install applies the chord route to an http router
func (t *Transporter) Install(server *Server, mux *mux.Router) {
	mux.Handle(t.findSuccessorPath, transporterHandler(t.findSuccessorHandler(server)))
	mux.Handle(t.findSuccessorPath, transporterHandler(t.notifyHandler(server)))
	mux.Handle(t.getPredecessorPath, transporterHandler(t.getPredecessorHandler(server)))
	mux.Handle(t.getSuccessorPath, transporterHandler(t.getSuccessorHandler(server)))
	mux.Handle(t.joinPath, transporterHandler(t.joinHandler(server))).Methods("POST")
	mux.Handle(t.startPath, transporterHandler(t.startHandler(server))).Methods("POST")
	mux.Handle(t.stopPath, transporterHandler(t.stopHandler(server))).Methods("POST")
	mux.Handle(t.getFingerTablePath, transporterHandler(t.getFingerTableHandler(server)))
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
		return nil, errors.Wrap(err, encodeReqErrMsg)
	}

	url := req.host + t.findSuccessorPath
	httpResp, err := t.httpClient.Post(url, "chord.protobuf", &b)
	if err != nil {
		return nil, errors.Wrap(err, "send successor request failed")
	}
	defer httpResp.Body.Close()

	successorResp := &FindSuccessorResponse{}
	if _, err = successorResp.Decode(httpResp.Body); err != nil {
		return nil, errors.Wrap(err, "send successor request failed")
	}

	return successorResp, nil
}

// SendNotifyRequest sends a request to other node to nofify it about the possible new predecessor
func (t *Transporter) SendNotifyRequest(server *Server, req *NotifyRequest) (*NotifyResponse, error) {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		return nil, errors.Wrap(err, encodeReqErrMsg)
	}

	url := req.targetHost + t.notifyPath
	httpResp, err := t.httpClient.Post(url, "chord.protobuf", &b)
	if err != nil {
		return nil, errors.Wrap(err, "send notify request failed")
	}
	defer httpResp.Body.Close()

	notifyResp := &NotifyResponse{}
	if _, err = notifyResp.Decode(httpResp.Body); err != nil {
		return nil, errors.Wrap(err, decodeRespErrMsg)
	}

	return notifyResp, nil
}

// SendGetPredecessorRequest sends a request to get the predecessor of server on given host
func (t *Transporter) SendGetPredecessorRequest(server *Server, host string) (*GetPredecessorResponse, error) {
	url := host + t.getPredecessorPath
	httpResp, err := t.httpClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "send getPredecessor request failed")
	}
	defer httpResp.Body.Close()

	predResp := &GetPredecessorResponse{}
	if _, err = predResp.Decode(httpResp.Body); err != nil {
		return nil, errors.Wrap(err, decodeRespErrMsg)
	}

	return predResp, nil
}

//	-------------------------------------------------------------------------
//
//	handler functions
//
//	-------------------------------------------------------------------------

const (
	encodeReqErrMsg  = "failed to encode request"
	decodeReqErrMsg  = "failed to decode request"
	encodeRespErrMsg = "faield to encode response"
	decodeRespErrMsg = "failed to decode response"
)

type transporterHandler func(http.ResponseWriter, *http.Request) error

func (h transporterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// findSuccessorHandler handles incoming request to find successor of the given key
func (t *Transporter) findSuccessorHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := &FindSuccessorRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			return errors.Wrap(err, decodeReqErrMsg)
		}

		resp, err := server.FindSuccessor(req)
		if resp == nil || err != nil {
			return errors.Wrap(err, "failed to find successor")
		}

		log.Printf("host %s's successor is: %s", req.host, resp.host)
		if _, err := resp.Encode(w); err != nil {
			return errors.Wrap(err, encodeRespErrMsg)
		}
		return nil
	}
}

// notifyHandler handles incoming notify about possibe new predecessor
func (t *Transporter) notifyHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := &NotifyRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			return errors.Wrap(err, decodeReqErrMsg)
		}

		//resp, err := server.processNotifyRequest(req)
		resp, err := server.notify(req)
		if resp == nil || err != nil {
			return errors.Wrap(err, "failed to notify")
		}

		if _, err := resp.Encode(w); err != nil {
			return errors.Wrap(err, encodeRespErrMsg)
		}
		return nil
	}
}

// getPredecessorHandler handles incoming request to return this local server's predecessor
func (t *Transporter) getPredecessorHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		predResp, err := server.processGetPredecessorRequest()
		if predResp == nil || err != nil {
			return errors.Wrap(err, "failed to process request of getting predecessor")
		}
		log.Printf("host %s's predecessor is %s", server.config.Host, predResp.host)

		if _, err := predResp.Encode(w); err != nil {
			return errors.Wrap(err, encodeRespErrMsg)
		}
		return nil
	}
}

// getSuccessorHandler handles the incoming request to return this node's successor
func (t *Transporter) getSuccessorHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		succResp, err := server.processGetSuccessorRequest()
		if err != nil {
			return errors.Wrap(err, "failed to process request of getting successor")
		}

		log.Printf("host %s's successor is %s", server.config.Host, succResp.host)

		if _, err := succResp.Encode(w); err != nil {
			return errors.Wrap(err, encodeRespErrMsg)
		}
		return nil
	}
}

// joinHandler handles the post request for this server to join an existing Chord ring
// the url pattern is '/join?host='
func (t *Transporter) joinHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		host := r.URL.Query().Get("host")
		return server.Join(host)
	}
}

// startHandler handles the incoming request to start this Chord server
func (t *Transporter) startHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return server.Start()
	}
}

// stopHandler handles incoming request to stop the Chord server
func (t *Transporter) stopHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return server.Stop()
	}
}

// getFingerTableHandler handles incoming request to log entries in finger table
func (t *Transporter) getFingerTableHandler(server *Server) transporterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fingers := server.node.Finger()
		for i := 0; i < server.config.HashBits; i++ {
			log.Printf("host %s's finger at index %d: %s", server.config.Host, i, fingers[i])
		}
		return nil
	}
}
