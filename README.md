# Chord
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

### Join Chord ring
By knowing the host name of another server that is participating in the Chord ring, this server can join the Chor ring as well
Suppose ```http://localhost:4000``` is a server in an existing Chord ring, we can join it by doing
```go
err := chordServer.Join("http://localhost:4000")
if err != nil {
    // handle error
}
```