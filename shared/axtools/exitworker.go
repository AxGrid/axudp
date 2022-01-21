package axtools

type Worker struct {
	Exit chan bool
}

func (w Worker) Stop() { w.Exit <- true }

func NewWorker() Worker {
	return Worker{
		Exit: make(chan bool, 1),
	}
}
