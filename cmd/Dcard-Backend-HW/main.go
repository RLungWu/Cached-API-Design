package main

import (
	"time"

	"github.com/RLungWu/Dcard-Backend-HW/api/ad"
	"github.com/RLungWu/Dcard-Backend-HW/api/admin"
	"github.com/RLungWu/Dcard-Backend-HW/api/public"

	"github.com/gin-gonic/gin"
)

const (
	MaxWorker = 10
	MaxQueue  = 20
)

// 一個可以發重工作請求的緩衝Channel
var AdRequestQueue chan ad.AdRequest

type Payload struct {
	AdRequest ad.AdRequest
}

type Job struct {
	Payload Payload
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

type Dispatcher struct {
	//Register th dispatcher's worker pool
	WorkerPool chan chan Job
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start function start a worker loop, listen quit channel
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				//Do the job
				//模擬工作時間
				time.Sleep(500 * time.Millisecond)
				//模擬工作完成
				println(job.Payload.AdRequest.Title)
			case <-w.quit:
				//Stop this worker
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quit <- true
	}()
}

func main() {
	d := NewDispatcher(MaxWorker)
	d.Run()

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/ad", admin.AdminCreateAD)
		v1.GET("/ad", public.GetAD)
	}

	router.Run(":8080")
}

func init() {
	AdRequestQueue = make(chan ad.AdRequest, MaxQueue)
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < MaxWorker; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case adReq := <-AdRequestQueue:
			go func(req ad.AdRequest) {
				
				job := Job{Payload: Payload{AdRequest: req}}
				//Try to catch an useable worker job channel
				jobChannel := <-d.WorkerPool

				//dispathch the job to the worker job channel
				jobChannel <- job
			}(adReq)
		}
	}
}
