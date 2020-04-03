package stages

import "time"

type b struct {
	stages []Stage
	prev   *Stage
}

func Builder() *b {
	return &b{prev: &Stage{Target: 1}}
}

func (b *b) RampUpAndSustain(target int, ramp, total time.Duration) *b {
	newStage := NewRampingStage(target, ramp, total)
	b.stages = append(b.stages, newStage)
	b.prev = &newStage
	return b
}

func (b *b) RampDown(ramp time.Duration) *b {
	newStage := NewRampDownStage(ramp, *b.prev)
	b.stages = append(b.stages, newStage)
	b.prev = &newStage
	return b
}

func (b *b) Build() []Stage {
	return b.stages
}
