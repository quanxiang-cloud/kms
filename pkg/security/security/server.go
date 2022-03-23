package security

import (
	"sync"
)

// Server server
type Server struct {
	lock sync.RWMutex

	commands commandTable
}

// New new
func New() *Server {
	return initServer()
}

func initServer() *Server {
	return &Server{
		commands: populateCommandTable(),
	}
}

func (s *Server) lookupCommand(cmd []Security, entity Type) (Ans, error) {
	if len(cmd) == 0 {
		return nil, nil
	}

	var err error
	var cp = entity

	for _, c := range cmd {

		s.lock.RLock()
		cp, err = s.getCmd(c.CMD).CommandProc(cp, c.Arg...)
		s.lock.RUnlock()

		if err != nil {
			return nil, err
		}
	}

	return Byte(cp)
}

func (s *Server) getCmd(name string) *CMD {
	cmd, ok := s.commands[name]
	if ok {
		return cmd
	}

	return s.commands["none"]
}

// Add 添加command
func (s *Server) Add(name string, command Command) {
	s.lock.RLock()
	copy := make(map[string]*CMD, len(s.commands)+1)
	for name, command := range s.commands {
		copy[name] = command
	}
	copy[name] = &CMD{Name: name, CommandProc: command}
	s.lock.RUnlock()

	s.change(copy)
}

// Remove 移除command
func (s *Server) Remove(name string) {
	s.lock.RLock()
	copy := make(map[string]*CMD, len(s.commands)-1)
	for key, command := range s.commands {
		if name == key {
			continue
		}
		copy[name] = command
	}
	s.lock.RUnlock()

	s.change(copy)
}

func (s *Server) change(commands map[string]*CMD) {
	s.lock.Lock()
	s.commands = commands
	s.lock.Unlock()
}
