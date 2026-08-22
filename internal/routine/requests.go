package routine

type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Weight      int    `json:"weight"`
	Coefficient int    `json:"coefficient"`
}

type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Weight      *int    `json:"weight"`
	Coefficient *int    `json:"coefficient"`
}
