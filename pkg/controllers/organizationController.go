package controllers

type NewOrganizationReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type InviteReq struct {
	Email string `json:"user_email" binding:"required"`
}

type UpdateOrganizationRes struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
