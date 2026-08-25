package plant

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
