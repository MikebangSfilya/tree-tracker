package plant

// ratioScale is a temporary mapping from integer DB values to the 0..1 JSON
// contract. Growth rules are out of scope; this is DTO conversion only.
const ratioScale = 100.0

// Response is the GET /api/plant JSON body (camelCase).
type Response struct {
	Epoch         int     `json:"epoch"`
	Phase         int     `json:"phase"`
	PhaseProgress float64 `json:"phaseProgress"`
	Branching     float64 `json:"branching"`
	Density       float64 `json:"density"`
	Curvature     float64 `json:"curvature"`
	Vitality      float64 `json:"vitality"`
	Seed          int64   `json:"seed"`
}

func NewResponse(p Plant) Response {
	return Response{
		Epoch:         p.Epoch,
		Phase:         p.Phase,
		PhaseProgress: toRatio(p.PhaseProgress),
		Branching:     toRatio(p.Branching),
		Density:       toRatio(p.Density),
		Curvature:     toRatio(p.Curvature),
		Vitality:      toRatio(p.Vitality),
		Seed:          p.Seed,
	}
}

func toRatio(v int) float64 {
	return float64(v) / ratioScale
}
