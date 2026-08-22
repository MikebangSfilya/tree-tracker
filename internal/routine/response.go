package routine

import "time"

type Response struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Weight      int       `json:"weight"`
	Coefficient int       `json:"coefficient"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewResponse(routine Routine) Response {
	return Response{
		ID:          routine.ID,
		Name:        routine.Name,
		Description: routine.Description,
		Category:    routine.Category,
		Weight:      routine.Weight,
		Coefficient: routine.Coefficient,
		CreatedAt:   routine.CreatedAt,
		UpdatedAt:   routine.UpdatedAt,
	}
}

func NewResponses(routines []Routine) []Response {
	responses := make([]Response, 0, len(routines))

	for _, routine := range routines {
		responses = append(responses, NewResponse(routine))
	}

	return responses
}
