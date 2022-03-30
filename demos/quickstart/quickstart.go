package main

import (
	"log"
	"time"

	"github.com/Xunzhuo/async"
)

func main() {
	q := async.Q().
		SetMaxWaitQueueLength(100).
		SetMaxWorkQueueLength(100).
		Start()
	q.AddJobAndRun(async.NewJob(longTimeJob).SetJobID("xunzhuo").SetSubID("a1"))
	q.AddJobAndRun(async.NewJob(longTimeJob).SetJobID("xunzhuo").SetSubID("a2"))
	q.AddJobAndRun(async.NewJob(longTimeJob).SetJobID("xunzhuo").SetSubID("a3"))
	q.AddJobAndRun(async.NewJob(longTimeJob).SetJobID("xunzhuo").SetSubID("a4"))
	q.AddJobAndRun(async.NewJob(longTimeJob).SetJobID("xunzhuo").SetSubID("a4"))
	q.AddJobAndRun(async.NewJob(longTimeJob).SetJobID("xunzhuo").SetSubID("a4"))

	q.AddJob(async.NewJob(longTimeJob).SetJobID("xunzhuo1").SetSubID("a4"))
	q.AddJob(async.NewJob(longTimeJob).SetJobID("xunzhuo1").SetSubID("a5"))
	q.AddJob(async.NewJob(longTimeJob).SetJobID("xunzhuo1").SetSubID("a5"))
	q.AddJob(async.NewJob(longTimeJob).SetJobID("xunzhuo1").SetSubID("a6"))

	j := q.GetJobByID("xunzhuo")
	log.Println(j.GetJobID(), j.GetSubIDs(), j.GetSubID())

	j = q.GetJobByID("xunzhuo1")
	log.Println(j.GetJobID(), j.GetSubIDs(), j.GetSubID())
}

func longTimeJob() {
	time.Sleep(500 * time.Millisecond)
}
