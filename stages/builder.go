package stages

import "time"

// B is a stages builder.
type B struct {
	stages     []Stage
	stagesToBe []StageToBe
}

// Builder returns a builder for creating a sequence of stages.
func Builder() *B {
	return &B{}
}

// RampUpAndSustain adds a ramp up and sustain stage to be, multiplying it by a potential multiplier.
// See NewRampingStage for more information.
func (b *B) RampUpAndSustain(target int, ramp, sustain time.Duration) *B {
	b.stagesToBe = append(b.stagesToBe, func(multiplier float64, _ Stage) Stage {
		return NewRampingStage(int(float64(target)*multiplier), ramp, sustain)
	})
	return b
}

// RampDown adds a ramp down to zero stage.
func (b *B) RampDown(ramp time.Duration) *B {
	b.stagesToBe = append(b.stagesToBe, func(multiplier float64, prev Stage) Stage {
		return NewRampDownStage(ramp, prev)
	})
	return b
}

// AddFixedStage adds a stage to be that will not be modified by a multiplier.
func (b *B) AddFixedStage(s Stage) *B {
	b.stagesToBe = append(b.stagesToBe, func(_ float64, _ Stage) Stage {
		return s
	})
	return b
}

// AddStageToBe adds a stage to be to the stagesToBe slice.
func (b *B) AddStageToBe(s StageToBe) *B {
	b.stagesToBe = append(b.stagesToBe, s)
	return b
}

// Build creates the stages from the stages-to-be in order using the given multiplier.
func (b *B) Build(multiplier float64) []Stage {
	for i, stageFunc := range b.stagesToBe {
		prev := Stage{Target: 1}
		if i-1 > 0 {
			prev = b.stages[i-1]
		}
		b.stages = append(b.stages, stageFunc(multiplier, prev))
	}
	return b.stages
}
