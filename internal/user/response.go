package user

type CreateResp struct {
	ID string `json:"id"`
}

type GetResp struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type DeleteResp struct {
	Deleted bool `json:"deleted"`
}
