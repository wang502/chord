# Chord
[![Go Report Card](https://goreportcard.com/badge/github.com/wang502/chord)](https://goreportcard.com/report/github.com/wang502/chord)
An Go implementation of Chord P2P Protocol 

## Install
Get the package
```
$ go get github.com/wang502/chord
```
Import the package
```go
import "github.com/wang502/chord"
```

## Usage
### Configure
Add a `config.json` file in your source folder
```json
{
  "Host": ,
  "HashBits": ,
  "NumNodes": ,
}
```

- ***Host***: the host name of ip of the local server that wants to join the Chord ring
- ***HashBits***: the number of bits in the hash bits to apply consistent hashing.
- ***NumNodes***: the max number of nodes to participate in Chord ring. `2^(HashBits) = NumNodes` 

Initialize config
- Initialize default congiguration with `HashBits=3 NumNodes=8` by passing only the host name
```go
host = "localhost:3000"
config := chord.DefaultConfig(host)
```

### Initialize Chord server
```go
import (
    "github.com/gorilla/mux"
	"github.com/wang502/chord"
)

transporter := chord.NewTransporter()
chordServer := chord.NewServer("chord1", config, transporter)
transporter.Install(server, mux.NewRouter())
```
Chord servers can communicate with each other using an HTTP transporter. And after transporter installs chord server, following url paths are mapped to respective handlers:
- "/findSuccessor": path to handle incoming request to find successor of an given id
- "/getPredecessor": path to return the predecessor of this chord node
- "/getSuccessor": path to return the successor of this chord node
- "/getFingerTable": path to return the finger table of this chord node
- "/notify": path to handle the notify request 
- "/join": path to handle a join request sent from a Chord server
- "/start": path to start this Chord server
- "/stop": path to stop this Chord server

### Join Chord ring
By knowing the host name of another server that is participating in the Chord ring, this server can join the Chor ring as well
This example joins a Chord ring consisting a ```http://localhost:4000``` host.
```go
err := chordServer.Join("http://localhost:4000")
if err != nil {
    // handle error
}
```

### Find successor
```go
succReq := NewFindSuccessorRequest(id, host)
succResp, err := chord.FindSuccessor(succReq)
```