package progress

import "github.com/MikebangSfilya/tree-tracker/internal/plant"

const pointsPerPhase = 100

func Advance(p plant.Plant, points int) plant.Plant {
	p.PhaseProgress += points
	for p.PhaseProgress >= pointsPerPhase {
		p.PhaseProgress -= pointsPerPhase
		if p.Phase < 3 {
			p.Phase++
		} else {
			p.Epoch++
		}
	}
	return p
}
