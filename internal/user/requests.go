package user


type CreateReq struct {
	Username string `json:"username"`
	Email	 string `json:"email"`
	Password string `json:"password"`
}

type GetAndDeleteReq struct {
	Email string `json:"email"`
}
