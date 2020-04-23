# github.com/astrogo-game/ticker

Description soon...

## Installation

```bash
$ go get -u -v github.com/astrogo-game/ticker
```

## Usage

```go
package main

import (

"github.com/astrago-game/ticker"
"time"
)

type myJob struct{
}
func (j *myJob) Tick(job *ticker.Job, now time.Time) {
    // Execute the logic
    
    // Mark the job as done. 
    job.Done()
    
    // Or Requeue
    job.Requeue()
    
    // Or MoveTo another worker
    
    job.MoveTo("another-worker")
}


func init() {    
	// Register an Worker
    df, _ := ticker.RegisterWorker("worker-name", ticker.NewGenericWorker("worker-name", time.Hour))
    
    job, err := ticker.RegisterJob("worker-name", &myJob{})
    
    // You can mark a Job as done before worker be execute and will be removed for the job list
    job.Done()

    // Start the ticker workers
    ticker.Start()

    // Stop the ticker workers
    ticker.Stop()

    // Execute the df func to remove/finish the worker
    df()
}
```
