package hwtimer

import (
	"fmt"
	"sync"
	"time"
)

// https://novoland.github.io/%E5%B9%B6%E5%8F%91/2014/07/26/%E5%AE%9A%E6%97%B6%E5%99%A8%EF%BC%88Timer%EF%BC%89%E7%9A%84%E5%AE%9E%E7%8E%B0.html
// https://medium.com/@raghavan99o/hashed-timing-wheel-2192b5ec8082
// https://www.cnkirito.moe/timer/
// http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html
// https://github.com/msackman/gotimerwheel/blob/master/gotimerwheel.go

// TODO
// 1. new task
// 2. run task
// 3. cancel task

type Timer struct {
	duration    time.Duration // tick duration, in milliseconds, default 100ms
	size        int64         // ticks per Wheel(Wheel Size), default 512
	wheel       []tasks
	wheelCursor int64
	m           sync.RWMutex
}

type Task struct {
	expiration time.Duration
	round      int64
	stopIndex  int64
	task       func()
}

type tasks map[*Task]bool

var TempChan = make(chan time.Time)

func NewTimer(duration time.Duration, size int64) *Timer {
	timer := &Timer{
		duration: duration * time.Millisecond,
		size:     size,
		wheel:    make([]tasks, size),
	}

	go func() {
		// TODO wheel cursor move
		for {
			s := time.Now().UTC()

			fmt.Println(timer.wheel)

			//tr := time.AfterFunc(timer.duration, func() {
			//	timer.wheelCursor = (timer.wheelCursor + 1) % timer.size
			//	fmt.Println(time.Now().UTC().Sub(s), timer.wheelCursor)
			//})

			c := time.NewTicker(100 * time.Millisecond)

			<-c.C

			timer.wheelCursor = (timer.wheelCursor + 1) % timer.size
			//fmt.Println(time.Now().UTC().Sub(s), timer.wheelCursor)

			go func(s time.Time) {
				timer.m.RLock()
				defer timer.m.RUnlock()

				//TempChan <- time.Now().UTC()
				//return

				tasks := timer.wheel[timer.wheelCursor]
				for task, _ := range tasks {
					task.round -= 1
					if task.round <= 0 {
						task.task()

						delete(tasks, task)
						fmt.Println("doing task: ", time.Now().UTC().Sub(s), timer.wheelCursor)
					}
				}

				fmt.Println(time.Now().UTC().Sub(s), tasks, timer.wheelCursor)

			}(s)

			//select {
			//case <-time.After(timer.duration):
			//	timer.wheelCursor = (timer.wheelCursor + 1) % timer.size
			//	//go func(s time.Time) {
			//	//	timer.m.RLock()
			//	//	defer timer.m.RUnlock()
			//	//
			//	//	//TempChan <- time.Now().UTC()
			//	//	//return
			//	//
			//	//	tasks := timer.wheel[timer.wheelCursor]
			//	//	for task, _ := range tasks {
			//	//		task.round -= 1
			//	//		if task.round <= 0 {
			//	//			task.task()
			//	//
			//	//			delete(tasks, task)
			//	//			fmt.Println("doing task: ", time.Now().UTC().Sub(s), timer.wheelCursor)
			//	//		}
			//	//	}
			//	//
			//	//	fmt.Println(time.Now().UTC().Sub(s), tasks, timer.wheelCursor)
			//	//
			//	//}(s)
			//
			//	fmt.Println(time.Now().UTC().Sub(s), timer.wheelCursor)
			//}
		}
	}()

	return timer
}

func (t *Timer) AfterFunc(d time.Duration, f func()) {
	t.schedule(&Task{
		expiration: d,
		task:       f,
	})
}

func (t *Timer) schedule(task *Task) {
	d := int64(task.expiration / t.duration)

	stopIndex := (t.wheelCursor + d%t.size) % t.size

	task.round = d / t.size
	task.stopIndex = stopIndex

	t.m.Lock()
	defer t.m.Unlock()

	ts := t.wheel[stopIndex]
	if ts == nil {
		ts = make(tasks)
	}
	ts[task] = true

	t.wheel[stopIndex] = ts

}

func (t *Timer) String() string {
	return fmt.Sprintf("%v\n", t.wheel)
}
