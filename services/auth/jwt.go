package auth

import (
	// "context"
	// "log"
	// "net/http"
	"strconv"
	"time"

	"github.com/capamir/go-api/configs"
	// "github.com/capamir/go-api/utils"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

// func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := utils.GetTokenFromRequest(r)

// 		token, err := validateJWT(tokenString)
// 		if err != nil {
// 			log.Printf("failed to validate token: %v", err)
// 			permissionDenied(w)
// 			return
// 		}

// 		if !token.Valid {
// 			log.Println("invalid token")
// 			permissionDenied(w)
// 			return
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		str := claims["userID"].(string)

// 		userID, err := strconv.Atoi(str)
// 		if err != nil {
// 			log.Printf("failed to convert userID to int: %v", err)
// 			permissionDenied(w)
// 			return
// 		}

// 		u, err := store.GetUserByID(userID)
// 		if err != nil {
// 			log.Printf("failed to get user by id: %v", err)
// 			permissionDenied(w)
// 			return
// 		}

// 		// Add the user to the context
// 		ctx := r.Context()
// 		ctx = context.WithValue(ctx, UserKey, u.ID)
// 		r = r.WithContext(ctx)

// 		// Call the function if the token is valid
// 		handlerFunc(w, r)
// 	}
// }