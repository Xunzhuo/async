package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Xunzhuo/async"
)

func main() {
	q := async.Q()

	q.SetMaxWaitQueueLength(10000).
		SetMaxWorkQueueLength(10000)
	q.Start()

	testMasterSlaveJob(q)

	// testRepeatMasterSlaveJob(q)

	q.Stop()
	time.Sleep(1 * time.Second)
}

func testMasterSlaveJob(q *async.Queue) {
	for i := 0; i < 1000; i++ {
		data := "a" + fmt.Sprint(i)
		job := async.NewJob(longTimeJob, data).SetJobID("xunzhuo").SetSubID(data)
		q.AddJobAndRun(job)
	}

	for {
		j1 := q.GetJobByID("xunzhuo")
		if datas, ok := q.GetJobsData(*j1); ok {
			for k, v := range datas {
				log.Println(k, v[0].(string))
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func testRepeatMasterSlaveJob(q *async.Queue) {
	for i := 0; i < 10; i++ {
		data := "a"
		job := async.NewJob(longTimeJob, data).SetJobID("xunzhuo").SetSubID(data)
		q.AddJobAndRun(job)
		// time.Sleep(1 * time.Second)
	}

	for {
		j := q.GetJobByID("xunzhuo")
		log.Println(j.GetJobID(), j.GetSubIDs(), j.GetSubID(), q.Length())
		if datas, ok := q.GetJobsData(*j); ok {
			for k, v := range datas {
				log.Println(k, v[0].(string))
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func testMultiSingleJob(q *async.Queue) {
	count := 1000
	for i := 0; i < count; i++ {
		data := "a" + fmt.Sprint(i)
		job := async.NewJob(longTimeJob, data).SetJobID(data)
		q.AddJobAndRun(job)
	}
	for i := 0; i < count; i++ {
		data := "a" + fmt.Sprint(i)
		m := q.GetJobByID(data)
		for {
			d, ok := q.GetJobData(*m)
			if ok {
				log.Println("SUCCESS", m.GetSubID(), d[0].(string))
				break
			}
			log.Println("MISSING", m.GetSubID(), m.GetJobID())
			time.Sleep(1 * time.Second)
		}
	}
}

func testRepeatSingleJob(q *async.Queue) {
	count := 1000
	for i := 0; i < count; i++ {
		data := "a"
		job := async.NewJob(longTimeJob, data).SetJobID(data)
		q.AddJobAndRun(job)
	}
	for i := 0; i < count; i++ {
		data := "a"
		m := q.GetJobByID(data)
		for {
			d, ok := q.GetJobData(*m)
			if ok {
				log.Println("SUCCESS", m.GetSubID(), d[0].(string))
				break
			}
			log.Println("MISSING", m.GetSubID(), m.GetJobID())
			time.Sleep(1 * time.Second)
		}
	}
}

func longTimeJob(name string) string {
	return "Hello World: " + name
}
