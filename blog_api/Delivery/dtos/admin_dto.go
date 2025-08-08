package dtos

// promote user request
type PromoteUserDTO struct {
	TargetUserID string `json:"target_user_id"`
}

// demote user request
type DemoteUserDTO struct {
	TargetUserID string `json:"target_user_id"`
}

// promote user response
type PromoteUserResponseDTO struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	NewRole string `json:"new_role"`
}

// demote user response
type DemoteUserResponseDTO struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	NewRole string `json:"new_role"`
}
