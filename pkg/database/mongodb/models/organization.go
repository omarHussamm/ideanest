package models

type OrganizationMember struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	AccessLevel string `json:"access_level" binding:"required"`
}

type Organization struct {
	ID                  uint                 `json:"id"`
	Name                string               `json:"name" binding:"required"`
	Description         string               `json:"description" binding:"required"`
	OrganizationMembers []OrganizationMember `json:"organization_members" binding:"required"`
}
