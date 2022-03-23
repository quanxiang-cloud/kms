package security

import "strings"

const (
	cmdSymbol   = "|"
	spaceSymBol = " "
)

func analysisString(cmd string) []Security {
	commands := strings.Split(cmd, cmdSymbol)
	ans := make([]Security, 0, len(commands))
	for _, command := range commands {
		params := strings.Split(command, spaceSymBol)

		size := len(params)
		if size == 0 {
			continue
		}

		cmd := Security{
			CMD: params[0],
		}
		if size > 1 {
			cmd.Arg = params[1:]
		}

		ans = append(ans, cmd)
	}

	return ans
}

// LookupCommandString command with string
func (s *Server) LookupCommandString(cmd string, entity Type) (Ans, error) {
	return s.lookupCommand(analysisString(cmd), entity)
}
