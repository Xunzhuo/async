package main

import (
	"log"
	"time"

	"github.com/Xunzhuo/async"
)

func main() {
	q := async.Q()

	q.SetMaxWaitQueueLength(100).
		SetMaxWorkQueueLength(100)
	q.Start()
	testSingleJob1(q)
	testMasterSlaveJob1(q)

	testSingleJob2(q)
	testMasterSlaveJob2(q)

	q.Stop()
}

func testMasterSlaveJob1(q *async.Queue) {
	q.AddJobAndRun(async.NewJob(longTimeJob, "a9").SetJobID("xunzhuo").SetSubID("a9"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a2").SetJobID("xunzhuo").SetSubID("a2"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a3").SetJobID("xunzhuo").SetSubID("a3"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a4").SetJobID("xunzhuo").SetSubID("a4"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a5").SetJobID("xunzhuo").SetSubID("a5"))

	j1 := q.GetJobByID("xunzhuo")
	log.Println(j1.GetJobID(), j1.GetSubIDs(), j1.GetSubID(), q.Length())
	for {
		if datas, ok := q.GetJobsData(*j1); ok {
			for k, v := range datas {
				log.Println(k, v[0].(string))
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func testMasterSlaveJob2(q *async.Queue) {
	q.AddJobAndRun(async.NewJob(longTimeJob, "a9").SetJobID("xunzhuo1").SetSubID("a9"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a2").SetJobID("xunzhuo1").SetSubID("a9"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a3").SetJobID("xunzhuo1").SetSubID("a9"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a4").SetJobID("xunzhuo1").SetSubID("a9"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "a5").SetJobID("xunzhuo1").SetSubID("a9"))

	j1 := q.GetJobByID("xunzhuo1")
	log.Println(j1.GetJobID(), j1.GetSubIDs(), j1.GetSubID(), q.Length())
	for {
		if datas, ok := q.GetJobsData(*j1); ok {
			for k, v := range datas {
				log.Println(k, v[0].(string))
			}
			break
		}
		time.Sleep(1 * time.Second)
	}

}

func testSingleJob1(q *async.Queue) {
	q.AddJobAndRun(async.NewJob(longTimeJob, "b1").SetJobID("xunzhuo11"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "b2").SetJobID("xunzhuo12"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "b3").SetJobID("xunzhuo13"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "b4").SetJobID("xunzhuo14"))

	m1 := q.GetJobByID("xunzhuo11")
	m2 := q.GetJobByID("xunzhuo12")
	m3 := q.GetJobByID("xunzhuo13")
	m4 := q.GetJobByID("xunzhuo14")
	for {
		d, ok := q.GetJobData(*m1)
		if ok {
			log.Println("SUCCESS", m1.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m1.GetSubID(), m1.GetJobID())
		time.Sleep(1 * time.Second)
	}
	for {
		d, ok := q.GetJobData(*m2)
		if ok {
			log.Println("SUCCESS", m2.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m2.GetSubID(), m2.GetJobID())
		time.Sleep(1 * time.Second)
	}
	for {
		d, ok := q.GetJobData(*m3)
		if ok {
			log.Println("SUCCESS", m3.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m3.GetSubID(), m3.GetJobID())
		time.Sleep(1 * time.Second)
	}
	for {
		d, ok := q.GetJobData(*m4)
		if ok {
			log.Println("SUCCESS", m4.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m4.GetSubID(), m4.GetJobID())
		time.Sleep(1 * time.Second)
	}
}

func testSingleJob2(q *async.Queue) {
	q.AddJobAndRun(async.NewJob(longTimeJob, "b1").SetJobID("xunzhuo21"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "b2").SetJobID("xunzhuo21"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "b3").SetJobID("xunzhuo21"))
	q.AddJobAndRun(async.NewJob(longTimeJob, "b4").SetJobID("xunzhuo21"))

	m1 := q.GetJobByID("xunzhuo21")
	m2 := q.GetJobByID("xunzhuo21")
	m3 := q.GetJobByID("xunzhuo21")
	m4 := q.GetJobByID("xunzhuo21")
	for {
		d, ok := q.GetJobData(*m1)
		if ok {
			log.Println("SUCCESS", m1.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m1.GetSubID(), m1.GetJobID())
		time.Sleep(1 * time.Second)
	}
	for {
		d, ok := q.GetJobData(*m2)
		if ok {
			log.Println("SUCCESS", m2.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m2.GetSubID(), m2.GetJobID())
		time.Sleep(1 * time.Second)
	}
	for {
		d, ok := q.GetJobData(*m3)
		if ok {
			log.Println("SUCCESS", m3.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m3.GetSubID(), m3.GetJobID())
		time.Sleep(1 * time.Second)
	}
	for {
		d, ok := q.GetJobData(*m4)
		if ok {
			log.Println("SUCCESS", m4.GetSubID(), d[0].(string))
			break
		}
		log.Println("MISSING", m4.GetSubID(), m4.GetJobID())
		time.Sleep(1 * time.Second)
	}
}

func longTimeJob(name string) string {
	return "Hello World: " + name
}
