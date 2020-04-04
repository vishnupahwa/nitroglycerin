package stages

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"reflect"
	"testing"
	"time"
)

func TestB_RampUpAndSustain(t *testing.T) {
	b := Builder()
	b.RampUpAndSustain(100, 1*time.Second, 2*time.Second)
	b.Build(0.5)
	got := b.stages
	want := []Stage{
		{
			Target:      50,
			StgDuration: 3 * time.Second,
			StgPacer: &SteadyUpPacer{
				upDuration: 1 * time.Second,
				min: vegeta.ConstantPacer{
					Freq: 1,
					Per:  time.Second,
				},
				max: vegeta.ConstantPacer{
					Freq: 50,
					Per:  time.Second,
				},
				slope:        0.000000000000000048999999999999995,
				minHitsPerNs: 0.000000001,
				maxHitsPerNs: 0.00000005,
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RampUpAndSustain() = %v, want %v", got, want)
	}
}
