package common

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true

	if err != nil {
		check = false
		return check, "109"
	}
	return check, "0"
}

func GenerateAllTokens(email string, userName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error, errcode string) {
	claims := &SignedDetails{
		Email:     email,
		Username:  userName,
		User_type: userType,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err, "103"
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string, collection *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})
	upsert := true
	filter := bson.M{"User_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	fmt.Println(userId)
	_, err := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}

func ErrorHandler(errcode string, response interface{}, totalcount int64, language string) HTTPResponse {
	if language == "" {
		language = "en"
	}
	ErrDetailes := Errorcodes[errcode].(map[string]interface{})
	msg := ErrDetailes[language].(string)
	Status := ErrDetailes["status"].(string)
	status, err := strconv.ParseInt(Status, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	code, err := strconv.ParseInt(errcode, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	return HTTPResponse{Msg: msg, Status: int(status), Data: response, TotalCount: totalcount, ErrorCode: code}
}

func CustomErrorHandler(errcode string, response interface{}, totalcount int64, language string, formater func(string) string) HTTPResponse {
	if language == "" {
		language = "en"
	}
	ErrDetailes := Errorcodes[errcode].(map[string]interface{})
	Msg := ErrDetailes[language].(string)
	msg := ""
	if formater != nil {
		msg = formater(Msg)
	}
	Status := ErrDetailes["status"].(string)
	status, err := strconv.ParseInt(Status, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	code, err := strconv.ParseInt(errcode, 10, 64)
	if err != nil {
		openlog.Error(err.Error())
		return HTTPResponse{Msg: "Internal server error", Status: 500}
	}
	return HTTPResponse{Msg: msg, Status: int(status), Data: response, TotalCount: totalcount, ErrorCode: code}
}
