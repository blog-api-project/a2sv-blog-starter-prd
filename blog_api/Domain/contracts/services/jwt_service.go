package services

type IJWTService interface {
	GenerateJWT(userID, role string) (string, error)
	ValidateJWT(tokenString string) (map[string]interface{}, error)
	GenerateRefreshToken(userID string , role string) (string, error)
	ValidateRefreshToken(tokenString string) (map[string]interface{}, error)
} 