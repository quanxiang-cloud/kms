package security

// Security security
type Security struct {
	CMD string
	Arg []string
}

// LookUpCommandSilce command with string
func (s *Server) LookUpCommandSilce(cmd []Security, entity []byte) (Ans, error) {
	return s.lookupCommand(cmd, entity)
}
