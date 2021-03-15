package authentication

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/macaron.v1"

	jwt "github.com/dgrijalva/jwt-go"
)

var MySigningKey = []byte("secret1234")

func JwtMiddleWare(ctx *macaron.Context) {

	fmt.Println("here authencation ")
	fmt.Println(ctx.Req.URL.Path)

	//if get then here, have given access to the user with or without token
	if ctx.Req.Method == "GET" || ctx.Req.URL.Path == "/register" || ctx.Req.URL.Path == "/request" || ctx.Req.URL.Path == "/edit-request" {
		log.Println(ctx.Req.Method, "request")
		ctx.Next()
	} else {
		//in case of other write operation, need to check if the user is valid
		fmt.Println("checking with jwt authencation")

		//retrieve the auth header
		authHeader := ctx.Req.Header.Get("Authorization")

		//if there is no auth header
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, "need auth token")
			return
		} else {
			// parsing the token from (bearer tokenString) with the secret key
			authHeaders := strings.Split(authHeader, " ")
			var token *jwt.Token
			var err error
			if len(authHeaders) < 2 {
				ctx.JSON(http.StatusUnauthorized, "no bearer token is provided")
				return
			}
			token, err = jwt.Parse(authHeaders[1], func(token *jwt.Token) (interface{}, error) {
				return MySigningKey, nil
			})
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, "no bearer token is provided")
				return
			}
			//checking if there is any claim and s the token valid
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println("user validity : ok jwt")

				//retrieving the claims info
				userId, okId := claims["userId"].(float64)

				//as i have set it capital i can use it outside the package
				CurrentUserId := int(userId)
				CurrentUserMail, okMail := claims["userMail"].(string)
				CurrentUserType, okType := claims["userType"].(string)

				if okId && okMail && okType {
					ctx.Req.Header.Set("current_user_id", string(CurrentUserId))
					ctx.Req.Header.Set("current_user_type", string(CurrentUserType))
					ctx.Req.Header.Set("current_user_mail", string(CurrentUserMail))
					ctx.Next()
				} else {
					ctx.JSON(http.StatusNotAcceptable, "the token is not valid .missing some info")
					return
				}

			} else {
				ctx.JSON(http.StatusUnauthorized, "need auth token")
				return
			}
		}
	}
}

func GenerateJWT(userMail string, userType string, userId int) (string, error) {
	//signing method declare
	token := jwt.New(jwt.SigningMethodHS256)

	//passing the parameter that i want to keep i my token
	claims := token.Claims.(jwt.MapClaims)
	claims["userMail"] = userMail
	claims["userType"] = userType
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 300).Unix()

	//generating the token string with my signing keys and claims parameter
	tokeString, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", err
	}
	fmt.Println("tokenString: ", tokeString)
	return tokeString, err

}

