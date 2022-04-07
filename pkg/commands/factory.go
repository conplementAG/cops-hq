package commands

func NewExecutor(logFileName string) *Executor {
	return &Executor{logFileName: logFileName}
}
