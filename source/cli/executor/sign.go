package executor

import (
	"fmt"
	"net"

	"../../trace"
	"../dispatcher"
)

type SignExecutor struct {
	dispatcher.Executor
}

func NewSignExecutor() *SignExecutor {
	return &SignExecutor{
		dispatcher.Executor{
			Name: "sign",
		},
	}
}

func (executor *SignExecutor) ShortHelp() string {
	return "short help for " + executor.Name + " command"
}

func (executor *SignExecutor) Validate() trace.ITrace {
	return nil
}

func (executor *SignExecutor) Execute() trace.ITrace {
	as, err := getMacAddr()
	if err != nil {
		fmt.Println(err)
	}
	for _, a := range as {
		fmt.Println(a)
	}

	return nil
}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
