// Copyright 2015 Donato Cassel Laynes Gonzales
//
// This file is part of GoOgame - Battle Simulator

package tools

import (
	"time"
)

type Task struct {
	Label     string
	StartTime time.Time
	EndTime   time.Time
	Position  int
}

var pos = 0

type Profiler struct {
	Tasks map[string]*Task
}

func (this *Profiler) Init(amount int) {
	this.Tasks = make(map[string]*Task, amount)
}

func (this *Profiler) StartTask(lbl string) {
	this.Tasks[lbl] = &Task{Label: lbl, StartTime: time.Now(), Position: pos}
	pos++
}

func (this *Profiler) EndTask(lbl string) {
	if this.Tasks[lbl] == nil {
		panic("Label " + lbl + " not found")
	}
	this.Tasks[lbl].EndTime = time.Now()
}

func (this *Profiler) GetTasks() []*Task {
	tasks := make([]*Task, len(this.Tasks))

	for _, tsk := range this.Tasks {
		tasks[tsk.Position] = tsk
	}

	return tasks
}
