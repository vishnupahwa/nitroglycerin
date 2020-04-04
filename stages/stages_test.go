package stages

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"reflect"
	"testing"
	"time"
)

func TestNewRampDownStage(t *testing.T) {
	type args struct {
		total time.Duration
		prev  Stage
	}
	tests := []struct {
		name string
		args args
		want Stage
	}{
		{
			name: "creates a stage which will ramp down to zero",
			args: args{
				total: 1 * time.Second,
				prev: Stage{
					Target:      100,
					StgDuration: 1 * time.Second,
					StgPacer:    vegeta.ConstantPacer{Freq: 100, Per: time.Second},
				},
			},
			want: Stage{
				Target:      0,
				StgDuration: 1 * time.Second,
				StgPacer: vegeta.LinearPacer{
					StartAt: vegeta.ConstantPacer{Freq: 100, Per: time.Second},
					Slope:   -100,
				},
			},
		},
		{
			name: "negative slope calculation is correct",
			args: args{
				total: 5 * time.Second,
				prev: Stage{
					Target:      100,
					StgDuration: 1 * time.Second,
					StgPacer:    vegeta.ConstantPacer{Freq: 100, Per: time.Second},
				},
			},
			want: Stage{
				Target:      0,
				StgDuration: 5 * time.Second,
				StgPacer: vegeta.LinearPacer{
					StartAt: vegeta.ConstantPacer{Freq: 100, Per: time.Second},
					Slope:   -20,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRampDownStage(tt.args.total, tt.args.prev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRampDownStage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRampingStage(t *testing.T) {
	type args struct {
		target  int
		ramp    time.Duration
		sustain time.Duration
	}
	tests := []struct {
		name string
		args args
		want Stage
	}{
		{
			name: "creates a stage which will ramp up to the target in the given time period and then sustain" +
				"at that target for the given time period",
			args: args{
				target:  100,
				ramp:    1 * time.Second,
				sustain: 2 * time.Second,
			},
			want: Stage{
				Target:      100,
				StgDuration: 3 * time.Second,
				StgPacer: &SteadyUpPacer{
					upDuration: 1 * time.Second,
					min: vegeta.ConstantPacer{
						Freq: 1,
						Per:  time.Second,
					},
					max: vegeta.ConstantPacer{
						Freq: 100,
						Per:  time.Second,
					},
					slope:        0.000000000000000099,
					minHitsPerNs: 0.000000001,
					maxHitsPerNs: 0.0000001,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRampingStage(tt.args.target, tt.args.ramp, tt.args.sustain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRampingStage() = %v, want %v", got, tt.want)
			}
		})
	}
}
