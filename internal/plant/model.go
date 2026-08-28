package plant

const (
	initialEpoch = 0
	initialPhase = 1
)

// Plant is the stored tree state. Visual parameters are integers in the database;
// the HTTP layer maps them to the float JSON contract.
type Plant struct {
	ID            int64
	Epoch         int
	Phase         int
	PhaseProgress int
	Branching     int
	Density       int
	Curvature     int
	Vitality      int
	Seed          int64
}

// InitialPlant is the starting tree from the issued contract: epoch 0, phase 1,
// zero visual progress. Seed is empty until an existing plant is reset.
func InitialPlant() Plant {
	return Plant{
		Epoch: initialEpoch,
		Phase: initialPhase,
	}
}
