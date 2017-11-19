package server

import (
	"time"
)

type Program struct {
	update func(delta float32)
	last   time.Time
}

func NewServerProgram(update func(delta float32)) Program {
	return Program{update, time.Now()}
}

func (program *Program) Update() {
	program.update(float32(time.Now().Sub(program.last).Nanoseconds()) / 1000000000.0)
	program.last = time.Now()
}
