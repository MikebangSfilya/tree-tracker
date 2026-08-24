package routine

type Response struct {
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	Weight      int      `json:"weight"`
	Coefficient int      `json:"coefficient"`
	TimeType    TimeType `json:"timeType"`
}

type TimeType struct {
	Temporary  bool `json:"temporary"`
	Disposable bool `json:"disposable"`
}

func NewResponse(routine Routine) Response {
	return Response{
		Category:    routine.Category,
		Type:        routine.Type,
		Weight:      routine.Weight,
		Coefficient: routine.Coefficient,
		TimeType: TimeType{
			Temporary:  routine.Temporary,
			Disposable: routine.Disposable,
		},
	}
}

func NewResponses(routines []Routine) []Response {
	responses := make([]Response, 0, len(routines))

	for _, routine := range routines {
		responses = append(responses, NewResponse(routine))
	}

	return responses
}
