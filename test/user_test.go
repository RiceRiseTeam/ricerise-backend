package test

import (
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
	"testing"
)

func TestUserRegister(test *testing.T) {
	//_, msg := postRequest[any]("/auth/register", &request.UserRegisterRequest{
	//	UserName: "testuser1",
	//	Email:    "181@outlook.com",
	//	NickName: "testnick",
	//	Password: "testpassword",
	//})
	//test.Log(msg)

	code, msg1 := postRequest[response.UserLoginResponse]("/auth/login", &request.UserLoginRequest{
		UserID:   "testuser1",
		Password: "testpassword",
	})
	test.Log(code, msg1)

	code, msg2 := postRequestWithAuth[any]("/user/logout", nil, msg1.AccessToken)

	test.Log(code, msg2)
}
