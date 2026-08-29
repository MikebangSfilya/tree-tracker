package progress

import (
	"testing"

	"github.com/MikebangSfilya/tree-tracker/internal/plant"
)

func TestAdvance(t *testing.T) {
	tests := []struct {
		name   string
		plant  plant.Plant
		points int
		want   plant.Plant
	}{
		{"adds points within phase", plant.Plant{PhaseProgress: 70}, 25, plant.Plant{PhaseProgress: 95}},
		{"moves to next phase", plant.Plant{Phase: 1, PhaseProgress: 90}, 25, plant.Plant{Phase: 2, PhaseProgress: 15}},
		{"moves mature tree to next epoch", plant.Plant{Epoch: 4, Phase: 3, PhaseProgress: 90}, 25, plant.Plant{Epoch: 5, Phase: 3, PhaseProgress: 15}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Advance(test.plant, test.points); got != test.want {
				t.Fatalf("Advance() = %#v, want %#v", got, test.want)
			}
		})
	}
}
