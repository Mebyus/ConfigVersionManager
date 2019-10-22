package command

import (
	"strings"
)

type Param struct {
	Name  string
	Value string
}

/**
Refactor Params using map
*/
type Command struct {
	Name   string
	Params []*Param
}

func (command *Command) Extract(name string) (value string, ok bool) {
	for _, param := range command.Params {
		if param.Name == name {
			value = param.Value
			ok = true
			return
		}
	}

	return
}

func parseParams(args []string, prefix string) []*Param {
	var param *Param

	params := make([]*Param, 0)
	for _, arg := range args {
		if strings.HasPrefix(arg, prefix) {
			if param != nil {
				params = append(params, param)
			}
			param = &Param{
				Name: arg[len(prefix):],
			}
		} else if param != nil {
			param.Value = arg
			params = append(params, param)
			param = nil
		} else {
			params = append(params, &Param{Value: arg})
		}
	}
	if param != nil {
		params = append(params, param)
	}

	return params
}

func ParseArgs(args []string, prefix string) *Command {
	if len(args) < 2 {
		commandPtr := &Command{
			Name:   "help",
			Params: nil,
		}
		return commandPtr
	}

	if strings.HasPrefix(args[1], prefix) {
		return nil
	}

	commandPtr := &Command{
		Name:   args[1],
		Params: parseParams(args[2:], prefix),
	}

	return commandPtr
}
