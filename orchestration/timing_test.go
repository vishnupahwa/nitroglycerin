package orchestration

import (
	"testing"
	"time"
)

func Test_calculateStartAt(t *testing.T) {
	type args struct {
		now     time.Time
		nextMin int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "finds nearest minute",
			args: args{
				now:     time.Unix(0, 0),
				nextMin: 1,
			},
			want: time.Unix(0, int64(1*time.Minute)).Unix(),
		},
		{
			name: "rounds to next 2nd minute",
			args: args{
				now:     time.Unix(0, 0),
				nextMin: 2,
			},
			want: time.Unix(0, int64(2*time.Minute)).Unix(),
		},
		{
			name: "rounds to next 2nd minute",
			args: args{
				now:     time.Unix(0, 0),
				nextMin: 2,
			},
			want: time.Unix(0, int64(2*time.Minute)).Unix(),
		},
		{
			name: "rounds to next 61st minute",
			args: args{
				now:     time.Unix(0, 0),
				nextMin: 61,
			},
			want: time.Unix(0, int64(61*time.Minute)).Unix(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateStartAt(tt.args.now, tt.args.nextMin); got != tt.want {
				t.Errorf("calculateStartAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
