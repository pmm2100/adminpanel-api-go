package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateJWT(ID int64, IsAdmin int) (string, error) {
	godotenv.Load()
	var secretKey = os.Getenv("JWT_SECRET")
	hours, _ := strconv.Atoi(os.Getenv("EXP_DATE"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":      ID,
		"IsAdmin": IsAdmin,
		"exp":     time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func CheckJWT(token string) (int64, int, error) {
	parsedtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unknown Signing method!")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, 0, errors.New("Could not parse token")
	}
	tokenValid := parsedtoken.Valid
	if !tokenValid {
		return 0, 0, errors.New("Token is not valid!")
	}
	claims, ok := parsedtoken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, errors.New("Invalid Claims!")
	}

	return int64(claims["ID"].(float64)), int(claims["IsAdmin"].(float64)), nil

}
