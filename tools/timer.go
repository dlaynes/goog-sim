package tools

import (
	"time"
)

type Task struct {
	StartTime time.Time
	EndTime   time.Time
}

type Profiler struct {
	Tasks map[string]*Task
}

func (this *Profiler) Init(amount int) {
	this.Tasks = make(map[string]*Task, amount)
}

func (this *Profiler) StartTask(lbl string) {
	this.Tasks[lbl] = &Task{StartTime: time.Now()}
}

func (this *Profiler) EndTask(lbl string) {
	if this.Tasks[lbl] == nil {
		panic("Label " + lbl + " not found")
	}
	this.Tasks[lbl].EndTime = time.Now()
}
