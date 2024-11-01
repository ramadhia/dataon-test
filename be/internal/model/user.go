package model

type Claim struct {
	ID       string   `json:"user_id" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Scope    []string `json:"scope" binding:"required"`
	Role     string   `json:"role"`
	IsClient bool     `json:"is_client"`
	Token    string
}

func (user Claim) IsClientToken() bool {
	return user.IsClient
}

func (user Claim) IsValidUserToken() bool {
	if user.IsClientToken() {
		return false
	}

	if hasPersonalId := user.ID != "" && user.Email != ""; !hasPersonalId {
		return false
	}

	return user.Token != "" && user.Role != ""
}
