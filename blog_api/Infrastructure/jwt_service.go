package infrastructure

import (
	services "blog_api/Domain/contracts/services"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTServiceImpl struct {
	secretKey []byte
}

func NewJWTService() services.IJWTService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key"
	}
	
	return &JWTServiceImpl{
		secretKey: []byte(secretKey),
	}
}

// generates a new JWT token
func (j *JWTServiceImpl) GenerateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // 15 minutes expiration
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//  validates a JWT token
func (j *JWTServiceImpl) ValidateJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token expired")
			}
		}
		if tokenType, ok := claims["type"].(string); ok {
			if tokenType != "access" {
				return nil, errors.New("invalid token type")
			}
		}

		return map[string]interface{}(claims), nil
	}

	return nil, errors.New("invalid token claims")
}

// generates a new refresh token
func (j *JWTServiceImpl) GenerateRefreshToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
        "role":    role,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
// validates a refresh token
func (j *JWTServiceImpl) ValidateRefreshToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token expired")
			}
		}

		if tokenType, ok := claims["type"].(string); ok {
			if tokenType != "refresh" {
				return nil, errors.New("invalid token type")
			}
		}

		return map[string]interface{}(claims), nil
	}

	return nil, errors.New("invalid token claims")
} 