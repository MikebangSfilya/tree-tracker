package user

import "time"

type CreateResp struct {
	ID string `json:"id"`
}

type GetResp struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
