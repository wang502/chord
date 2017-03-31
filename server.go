package chord

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// -------------------------------------------------------------------------
//
// Constants
//
// -------------------------------------------------------------------------

const (
	// DefaultStabilizeInterval is the interval that this server will start the stabilize process
	DefaultStabilizeInterval = 50 * time.Millisecond

	// DefaultFixFingerInterval is the interval that this server will repeat fixing its finger table
	DefaultFixFingerInterval = 50 * time.Millisecond
)

const (
	// Stopped denotes that Chord server is not started yet or has been stooped
	Stopped = "stopped"

	// Running denotes that Chord server is currently running
	Running = "running"
)

type event struct {
	value interface{}
	res   interface{}
	c     chan error
}

// Server represents a single node in Chord protocol
type Server struct {
	name  string
	state string
	node  *Node
	sync.RWMutex
	config      *Config
	transporter *Transporter

	stabilizeInterval time.Duration
	fixFingerInterval time.Duration

	stopChan chan bool

	routineGroup sync.WaitGroup

	c chan *event
}

// NewServer initializes a new local server involved in Chord protocol
func NewServer(name string, config *Config, transporter *Transporter) *Server {
	server := &Server{
		name:              name,
		state:             Stopped,
		node:              NewNode(config),
		config:            config,
		transporter:       transporter,
		stabilizeInterval: DefaultStabilizeInterval,
		fixFingerInterval: DefaultFixFingerInterval,
		stopChan:          make(chan bool),
		c:                 make(chan *event, 200),
	}
	return server
}

// Do tries to execute the command and returns the result
func (server *Server) Do(command interface{}) (interface{}, error) {
	return server.sendCommand(command)
}

// Join joins an existing chord ring, given existingHost is one of the node in the ring
func (server *Server) Join(existingHost string) error {
	localNode := server.node
	findSuccessorReq := NewFindSuccessorRequest(localNode.ID, existingHost)
	findSuccessorResp, err := server.transporter.SendFindSuccessorRequest(server, findSuccessorReq)
	if err != nil {
		return fmt.Errorf("Chord.join.error.%s", err)
	}

	successorNode := NewRemoteNode([]byte(findSuccessorResp.ID), findSuccessorResp.host)
	localNode.SetSuccessor(successorNode)

	if err = server.stabilize(); err != nil {
		return fmt.Errorf("Chord.join.error.%s", err)
	}
	log.Printf("[Join]host %s joined Chord ring", server.config.Host)
	return nil
}

// Leave the chord ring
func (server *Server) Leave() error {
	return nil
}

// Start the Chord server
func (server *Server) Start() error {
	if server.Running() {
		return fmt.Errorf("chord.Start.error.%s", server.state)
	}

	server.stopChan = make(chan bool)
	server.SetState(Running)

	server.routineGroup.Add(1)
	go func() {
		defer server.routineGroup.Done()
		server.startPeriodicalStabilize()
	}()

	server.routineGroup.Add(1)
	go func() {
		defer server.routineGroup.Done()
		server.startPeriodicalFixFinger()
	}()

	server.routineGroup.Add(1)
	go func() {
		defer server.routineGroup.Done()
		server.eventLoop()
	}()

	return nil
}

// Stop the Chord server
func (server *Server) Stop() error {
	log.Printf("stopping Chord server %s......", server.config.Host)
	if server.State() == Stopped {
		return fmt.Errorf("chord.Stop.error.%s", server.State())
	}

	close(server.stopChan)

	// make sure all goroutines are stopped
	server.routineGroup.Wait()
	server.SetState(Stopped)
	log.Printf("stopped Chord server %s", server.config.Host)
	return nil
}

// Running checks if Chord server is running
func (server *Server) Running() bool {
	server.Lock()
	defer server.Unlock()
	return server.state == Running
}

// sendCommand sends command to be executed into command channel and block waiting for result
func (server *Server) sendCommand(command interface{}) (interface{}, error) {
	if !server.Running() {
		return nil, errors.New("chord.sendCommand.error:server is not running")
	}
	e := &event{
		value: command,
		c:     make(chan error),
	}

	select {
	case server.c <- e:
	case <-server.stopChan:
		return nil, errors.New("chord.sendCommand.error: Server Stopped")
	default:
	}

	select {
	case <-server.stopChan:
		return nil, errors.New("chord.sendCommand.error: Server Stopped")
	case err := <-e.c:
		return e.res, err
	}
}

// eventLoop handles the incoming commands to be executed, notify request
func (server *Server) eventLoop() {
	var err error

	stopChan := server.stopChan
	state := server.State()
	for state != Stopped {
		select {
		case <-stopChan:
			log.Printf("chord.PeriodicalFixFinger.stop.%s", server.config.Host)
			return
		case ev := <-server.c:
			switch req := ev.value.(type) {
			case Command:
				ev.res, err = server.processCommand(req)
			case *NotifyRequest:
				ev.res, err = server.processNotifyRequest(req)
			default:
				err = errors.New("Command did not implements Apply() method")
			}
			ev.c <- err
		}
		state = server.State()
	}
}

// process the command to be applied on this Chord server
func (server *Server) processCommand(c Command) (interface{}, error) {
	return c.Apply(server)
}

// -------------------------------------------------------------------------
//
// Finger table
//
// -------------------------------------------------------------------------

// startPeriodicalFixFinger starts the periodical process of fixing finger table
func (server *Server) startPeriodicalFixFinger() {
	c := make(chan bool)
	server.routineGroup.Add(1)
	go func() {
		defer server.routineGroup.Done()
		server.periodicalFixFinger(c)
	}()
	<-c
}

func (server *Server) periodicalFixFinger(c chan bool) {
	c <- true

	stopChan := server.stopChan
	ticker := time.Tick(server.fixFingerInterval)

	log.Printf("chord.PeriodicalFixFinger.host: %s.interval: %s", server.config.Host, server.fixFingerInterval)

	state := server.State()
	for state != Stopped {
		select {
		case <-stopChan:
			log.Printf("chord.PeriodicalFixFinger.stop.%s", server.config.Host)
			return
		case <-ticker:
			err := server.fixFinger()
			if err != nil {
				log.Printf("[ERROR]%s.chord.PeriodicalFixFinger.error.%s", server.config.Host, err)
			}
		}

		state = server.State()
	}
}

func (server *Server) fixFinger() error {
	hb := server.config.HashBits

	node := server.node
	finger := node.Finger()
	node.fingerIndex = node.fingerIndex + 1
	next := node.fingerIndex
	if next >= hb {
		node.fingerIndex = 0
		next = 0
	}

	if finger[next] == nil {
		finger[next] = &FingerEntry{
			start: powerOffset(node.ID, next, hb),
		}
	}
	succReq := NewFindSuccessorRequest(finger[next].start, "")
	succResp, err := server.FindSuccessor(succReq)
	if err != nil {
		return err
	}

	finger[next].node = []byte(succResp.ID)
	finger[next].host = succResp.host

	//log.Printf("[DEBUG]%s's successor is %s", server.config.Host, server.node.Successor().host)
	log.Printf("[Fix Finger]%s's finger entry at %d is %s", server.config.Host, next, succResp.host)

	return nil
}

// startPeriodicalStabilize starts start the periodical stabilizing process
func (server *Server) startPeriodicalStabilize() {
	c := make(chan bool)
	server.routineGroup.Add(1)
	go func() {
		defer server.routineGroup.Done()
		server.periodicalStabilize(c)
	}()
	<-c
}

func (server *Server) periodicalStabilize(c chan bool) {
	c <- true

	stopChan := server.stopChan
	ticker := time.Tick(server.stabilizeInterval)

	log.Printf("chord.PeriodicalStabilize.host: %s.interval: %s", server.config.Host, server.stabilizeInterval)

	state := server.State()

	for state != Stopped {
		select {
		case <-stopChan:
			log.Printf("chord.PeriodicalStabilize.stop.%s", server.config.Host)
			return
		case <-ticker:
			err := server.stabilize()
			if err != nil {
				log.Printf("[ERROR]%s.chord.PeriodicalStabilize.error.%s", server.config.Host, err)
			}
		}

		state = server.State()
	}
}

// stabilize is called periodically to verify this server's immediate successor and tells the successor about this server
func (server *Server) stabilize() error {
	if server.node.Successor() == nil {
		return fmt.Errorf("no need to stabilize.no successor")
	}

	successor := server.node.Successor()
	predResp, err := server.transporter.SendGetPredecessorRequest(server, successor.host)
	if predResp == nil {
		log.Printf("[ERROR]stabilize.error.%s", err)
	} else if err != nil {
		return fmt.Errorf("stabilize.error.%s", err)
	} else {
		ID := []byte(predResp.ID)
		host := predResp.host

		if server.config.Host == successor.host {
			// if this node is same as its successor, then we update the successor to be the predecessor,
			// since there are at most 2 nodes in the ring now
			server.node.SetSuccessor(NewRemoteNode(ID, host))

		} else if between(server.node.ID, successor.ID, ID) {
			// verifies server's immediate successor
			// if the successor's predecessor has an ID bigger than this server, then it means this server's immediate successor
			// should be updated to the one contained in the response

			server.node.SetSuccessor(NewRemoteNode(ID, host))

		}
	}

	// notify the immediate successor about the server
	_, err = server.transporter.SendNotifyRequest(server, NewNotifyRequest(server.node.ID, server.config.Host, server.node.Successor().host))
	if err != nil {
		return fmt.Errorf("stabilize.error.%s", err)
	}
	log.Printf("[Stabilize]%s(%x)'s successor is %s(%x)", server.config.Host, server.node.ID, server.node.Successor().host, server.node.Successor().ID)

	return nil
}

func (server *Server) notify(req *NotifyRequest) (*NotifyResponse, error) {
	res, err := server.sendCommand(req)
	if res != nil {
		return res.(*NotifyResponse), err
	}
	return &NotifyResponse{}, err
}

// Notify handles the NotifyRequest sent from another server
func (server *Server) processNotifyRequest(req *NotifyRequest) (*NotifyResponse, error) {
	possiblePredID := []byte(req.ID)
	possiblePredHost := req.host
	currentPredecessor := server.node.Predecessor()

	// when this node haven't set its predecessor, then new incoming notify request is from a node that should be a predecessor
	if currentPredecessor == nil {
		server.node.SetPredecessor(NewRemoteNode(possiblePredID, possiblePredHost))
		return NewNotifyResponse(server.node.ID, server.config.Host), nil
	}
	// update the predecessor if the notify request is from a node that has bigger byte value than the current predecessor
	if between(currentPredecessor.ID, server.node.ID, possiblePredID) {
		server.node.SetPredecessor(NewRemoteNode(possiblePredID, possiblePredHost))
		return NewNotifyResponse(server.node.ID, server.config.Host), nil
	}
	return &NotifyResponse{}, nil
}

// FindSuccessor handles a incoming request sent from other server to help find successor
func (server *Server) FindSuccessor(req *FindSuccessorRequest) (*FindSuccessorResponse, error) {
	id := []byte(req.ID)
	localNode := server.node
	resp := &FindSuccessorResponse{}

	successor := localNode.Successor()
	if betweenRightIncl(localNode.ID, successor.ID, id) {
		resp.ID = string(successor.ID)
		resp.host = successor.host
		return resp, nil
	}

	closestPre := server.closestPreceedingNode(id)
	if closestPre == nil {
		resp.ID = string(localNode.ID)
		resp.host = server.config.Host
		return resp, nil
	}
	findSuccRequest := NewFindSuccessorRequest(id, closestPre.host)
	return server.transporter.SendFindSuccessorRequest(server, findSuccRequest)
}

// closestPreceedingNode is a helper function to find the cloest preceding node of the node with given hashed id from finger table
func (server *Server) closestPreceedingNode(id []byte) *RemoteNode {
	localNode := server.node
	finger := localNode.Finger()
	for i := server.config.HashBits - 1; i >= 0; i-- {
		if finger[i] != nil {
			if finger[i].node != nil && finger[i].host != "" {
				if between(localNode.ID, id, finger[i].node) {
					return &RemoteNode{ID: finger[i].node, host: finger[i].host}
				}
			}
		}
	}
	return nil
}

// HandleGetPredecessorRequest returns the predecessor of this local node
func (server *Server) processGetPredecessorRequest() (*GetPredecessorResponse, error) {
	pred := server.node.Predecessor()
	if pred == nil {
		return nil, fmt.Errorf("this node has no predecessor")
	}
	resp := NewGetPredecessorResponse(pred.ID, pred.host)
	return resp, nil
}

// HandleGetSuccessorRequest returns the successor of this local node
func (server *Server) processGetSuccessorRequest() (*FindSuccessorResponse, error) {
	succ := server.node.Successor()
	resp := &FindSuccessorResponse{}
	resp.ID = string(succ.ID)
	resp.host = succ.host

	return resp, nil
}

// -------------------------------------------------------------------------
//
// Getter
//
// -------------------------------------------------------------------------

// State retrieves the current state of Chord server
func (server *Server) State() string {
	server.Lock()
	defer server.Unlock()
	return server.state
}

// -------------------------------------------------------------------------
//
// Setter
//
// -------------------------------------------------------------------------

// SetStabilizeInterval sets the interval of periodical stabilizing
func (server *Server) SetStabilizeInterval(duration time.Duration) {
	server.Lock()
	defer server.Unlock()
	server.stabilizeInterval = duration
}

// SetFixFingerInterval sets the interval of periodical process of fixing finger table
func (server *Server) SetFixFingerInterval(duration time.Duration) {
	server.Lock()
	defer server.Unlock()
	server.fixFingerInterval = duration
}

// SetState sets the current state of Chord server
func (server *Server) SetState(state string) {
	server.Lock()
	defer server.Unlock()
	server.state = state
}
