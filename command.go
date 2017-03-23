package chord

// Command represents an interface to apply a command on a Chord server
type Command interface {
	CommandName() string
	Apply(s *Server) (interface{}, error)
}
