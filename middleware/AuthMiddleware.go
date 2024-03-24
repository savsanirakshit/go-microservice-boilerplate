package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

var secretKey = []byte("secret-key")

func AuthMiddleware(next http.Handler) http.Handler {

	// Todo : User this for http2
	//
	//if privateKey == nil || publicKey == nil {
	//	privateKey = LoadPrivateKey(common.CurrentWorkingDir() + "/server.key")
	//	publicKey = LoadPublicKey(common.CurrentWorkingDir() + "/server-pub.key")
	//}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !(strings.Contains(r.URL.Path, "/login") || strings.Contains(r.URL.Path, "/home")) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Missing authorization header")
				return
			}
			tokenString = tokenString[len("Bearer "):]

			err := verifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Invalid token")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// TODO : User this for http2
//
//func LoadPrivateKey(path string) *rsa.PrivateKey {
//	keyBytes, err := ioutil.ReadFile(path)
//	if err != nil {
//		log.Fatalf("Error reading private key: %v", err)
//	}
//
//	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
//	if err != nil {
//		log.Fatalf("Error parsing private key: %v", err)
//	}
//
//	return key
//}
//
//func LoadPublicKey(path string) *rsa.PublicKey {
//	keyBytes, err := ioutil.ReadFile(path)
//	if err != nil {
//		log.Fatalf("Error reading public key: %v", err)
//	}
//
//	key, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
//	if err != nil {
//		log.Fatalf("Error parsing public key: %v", err)
//	}
//
//	return key
//}

func verifyToken(tokenString string) error {
	// Todo : for http2
	//
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	// Validate the algorithm
	//	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	//	}
	//	// Return the public key for verification
	//	return publicKey, nil
	//})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = 1
	claims["username"] = "Rakshit"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours
	// Todo : use for http2
	//
	//tokenString, err := token.SignedString(privateKey)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
