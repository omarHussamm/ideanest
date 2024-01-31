package controllers

type LoginReq struct {
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type SignUpReq struct {
	Name     string `json:"name,omitempty" binding:"required"`
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type TokenReqBody struct {
	Refresh_token string `json:"refresh_token,omitempty" binding:"required"`
}
