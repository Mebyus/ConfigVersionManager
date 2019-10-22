package dispatcher

import (
	"fmt"

	"../../log"
	"../../trace"
	"../command"
	"../executor"
)

type Dispatcher struct {
	pool          []executor.IExecutor
	etraceFactory trace.IErrorTraceFactory
	logger        log.ILogger
}

func NewDispatcher(logger log.ILogger, factory trace.IErrorTraceFactory) *Dispatcher {
	return &Dispatcher{
		pool:          make([]executor.IExecutor, 0),
		logger:        logger,
		etraceFactory: factory,
	}
}

/**
Registers an executor in the dispatcher. So it becomes available for search in
Dispatch method.
*/
func (dispatcher *Dispatcher) Register(executor executor.IExecutor) {
	executor.ChangeEtraceFactory(dispatcher.etraceFactory)
	executor.ChangeLogger(dispatcher.logger)
	dispatcher.pool = append(dispatcher.pool, executor)
}

func (dispatcher *Dispatcher) ChangeEtraceFactory(etraceFactory trace.IErrorTraceFactory) {
	dispatcher.etraceFactory = etraceFactory
	for _, executor := range dispatcher.pool {
		executor.ChangeEtraceFactory(etraceFactory)
	}
}

func (dispatcher *Dispatcher) ChangeLogger(logger log.ILogger) {
	dispatcher.logger = logger
	for _, executor := range dispatcher.pool {
		executor.ChangeLogger(logger)
	}
}

func (dispatcher *Dispatcher) Dispatch(command *command.Command) {
	if command == nil {
		err := fmt.Errorf("invalid command format: command can't start with prefix")
		etrace := dispatcher.etraceFactory.CreateTrace(err, log.ERROR)
		dispatcher.logger.LogTrace(etrace)
		return
	}

	executor := dispatcher.search(command.Name)
	if executor == nil {
		err := fmt.Errorf("unrecognized command \"%s\"", command.Name)
		etrace := dispatcher.etraceFactory.CreateTrace(err, log.ERROR)
		dispatcher.logger.LogTrace(etrace)
		return
	}

	executor.Load(command)

	etrace := executor.Validate()
	if etrace != nil {
		etrace.Add(fmt.Sprintf("Preparing \"%s\" for execution", command.Name))
		dispatcher.logger.LogTrace(etrace)
		return
	}

	etrace = executor.Execute()
	if etrace != nil {
		etrace.Add(fmt.Sprintf("Something went wrong while executing \"%s\" command", command.Name))
		dispatcher.logger.LogTrace(etrace)
	}

	return
}

func (dispatcher *Dispatcher) search(name string) executor.IExecutor {
	for _, executor := range dispatcher.pool {
		if executor.Match(name) {
			return executor
		}
	}

	return nil
}
