package test

import (
	"ricerise/internal/dto/request"
	"testing"
)

func TestUserRegister(test *testing.T) {
	_, msg := postRequest("/user/register", &request.UserRegisterRequest{})
	test.Log(msg)
}
