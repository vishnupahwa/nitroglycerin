package stages

import "time"

type b struct {
	stages []Stage
	prev   *Stage
}

func Builder() *b {
	return &b{prev: &Stage{Target: 1}}
}

func (b *b) Add(target int, duration time.Duration) *b {
	newStage := NewStage(target, duration, *b.prev)
	b.stages = append(b.stages, newStage)
	b.prev = &newStage
	return b
}

func (b *b) AddWithModifiers(target int, duration time.Duration, before func() error, after func() error) *b {
	newStage := NewStage(target, duration, *b.prev)
	if before() != nil {
		newStage.Before = before
	}
	if after() != nil {
		newStage.After = after
	}
	b.stages = append(b.stages, newStage)
	b.prev = &newStage
	return b
}

func (b *b) Build() []Stage {
	return b.stages
}
