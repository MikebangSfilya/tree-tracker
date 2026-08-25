package routine

type CreateRequest struct {
	Category    string `json:"category"`
	Type        string `json:"type"`
	Weight      int    `json:"weight"`
	Coefficient int    `json:"coefficient"`
	Temporary   bool   `json:"temporary"`
	Disposable  bool   `json:"disposable"`
}

type UpdateRequest struct {
	Category    *string `json:"category"`
	Type        *string `json:"type"`
	Weight      *int    `json:"weight"`
	Coefficient *int    `json:"coefficient"`
	Temporary   *bool   `json:"temporary"`
	Disposable  *bool   `json:"disposable"`
}
