package test

import (
	"fmt"
	"ricerise/internal/dto"
	"ricerise/internal/dto/request"
	"ricerise/internal/dto/response"
	"testing"
)

func TestUserRegister(test *testing.T) {
	code, msg := postRequest[any]("/auth/register", &request.UserRegisterRequest{
		UserName: "testuser2",
		Email:    "181@outlook.com",
		NickName: "testnick",
		Password: "testpassword",
	})
	test.Log(code, msg)

	code, msg1 := postRequest[response.UserLoginResponse]("/auth/login", &request.UserLoginRequest{
		UserID:   "testuser2",
		Password: "testpassword",
	})
	test.Log(code, msg1)

	code, msg2 := postRequestWithAuth[any]("/user/logout", nil, msg1.AccessToken)

	test.Log(code, msg2)
}

func TestAdminOperation(test *testing.T) {
	code, msg1 := postRequest[response.UserLoginResponse]("/auth/login", &request.UserLoginRequest{
		UserID:   "testuser2",
		Password: "testpassword",
	})
	test.Log(code, msg1)

	code, msg2 := getRequestWithAuth[response.AdminCommentReviewList]("/admin/comments", nil, msg1.AccessToken)
	test.Log(code, msg2)

	code, msg3 := postRequestWithAuth[dto.LocationDto]("/map/locations", &request.UploadLocationRequest{
		Name:        "测试地点1",
		Latitude:    1.0,
		Longitude:   2.0,
		Description: "这是测试地点1",
	}, msg1.AccessToken)

	test.Log(code, msg3)

	code, msg4 := postRequestWithAuth[[]*dto.LocationDto]("/map/location/view", &request.GetLocationsRequest{
		MinLng: 0.1012,
		MinLat: 0.1001,
		MaxLat: 10.102,
		MaxLng: 10.123,
	}, msg1.AccessToken)

	test.Log(code, (*msg4)[1])

	code, msg5 := postRequestWithAuth[any](fmt.Sprintf("/map/location/%d/comments", (*msg4)[1].ID), &request.UploadCommentRequest{Rating: 5, Content: "牛逼"}, msg1.AccessToken)

	test.Log(code, msg5)

	code, msg6 := getRequestWithAuth[[]*dto.CommentDto](fmt.Sprintf("/map/location/%d/comments", (*msg4)[1].ID), nil, msg1.AccessToken)
	test.Log(code, msg6)
}
