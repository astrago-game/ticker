package ticker

func Start() {
	ticker.Start()
}

func Stop() {
	ticker.Stop()
}

func RegisterWorker(name string, worker Worker) (DeleteFn, error) {
	return ticker.RegisterWorker(name, worker)
}

func RegisterJob(workerName string, object JobObject) (*Job, error) {
	return ticker.RegisterJob(workerName, object)
}

func RequeueJob(job *Job) error {
	return ticker.RequeueJob(job)
}

func GetManager() Manager {
	return ticker
}
