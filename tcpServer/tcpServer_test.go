package main

import (
	"testing"
	"usermana/protocol"
)

var token string

// TestSignUpService 测试用户注册函数SignUpService.
func TestSignUpService(t *testing.T) {
	var tests = []struct {
		req protocol.ReqSignUp
		ret int
	}{
		{protocol.ReqSignUp{
			UserName: "botSignUp1",
			Password: "123",
			NickName: "botAABB",
		}, 0},
	}
	for _, test := range tests {
		resp := SignUpService(test.req)
		if resp.Ret != test.ret {
			t.Errorf("SignUpService didn't pass. username:%s, password:%s, nickname:%s, ret:%d", test.req.UserName, test.req.Password, test.req.NickName, test.ret)
		}
	}
}

//TestLoginService 测试用的登陆函数LoginService.
func TestLoginService(t *testing.T) {
	var tests = []struct {
		req protocol.ReqLogin
		ret int
	}{
		{protocol.ReqLogin{
			UserName: "botSignUp1",
			Password: "123",
		}, 0},
	}
	for _, test := range tests {
		resp := LoginService(test.req)
		if resp.Ret != test.ret {
			t.Errorf("LoginService didn't pass. username:%s, password:%s, ret:%d", test.req.UserName, test.req.Password, test.ret)
		} else {
			token = resp.Token
		}
	}
}

// TestGetProfileService 测试获取用户信息函数TestGetProfileService.
func TestGetProfileService(t *testing.T) {
	var tests = []struct {
		req protocol.ReqGetProfile
		ret int
	}{
		{protocol.ReqGetProfile{
			UserName: "botSignUp1",
			Token:    token,
		}, 0},
	}
	for _, test := range tests {
		resp := GetProfileService(test.req)
		if resp.Ret != test.ret {
			t.Errorf("GetProfileService didn't pass. username:%s, ret:%d", test.req.UserName, test.ret)
		}
	}
}

// TestUpdateProfilePicService 测试更新用户信息函数UpdateProfilePicService.
func TestUpdateProfilePicService(t *testing.T) {
	var tests = []struct {
		req protocol.ReqUpdateProfilePic
		ret int
	}{
		{protocol.ReqUpdateProfilePic{
			UserName: "botSignUp1",
			FileName: "http://127.0.0.1:1188/static/default.jpeg",
			Token:    token,
		}, 0},
	}
	for _, test := range tests {
		resp := UpdateProfilePicService(test.req)
		if resp.Ret != test.ret {
			t.Errorf("UpdateProfilePicService didn't pass. username:%s, filename:%s, ret:%d", test.req.UserName, test.req.FileName, test.ret)
		}
	}
}

// TestUpdateNickNameService 测试更新用户昵称函数UpdateNickNameService.
func TestUpdateNickNameService(t *testing.T) {
	var tests = []struct {
		req protocol.ReqUpdateNickName
		ret int
	}{
		{protocol.ReqUpdateNickName{
			UserName: "botSignUp1",
			NickName: "bot1188",
			Token:    token,
		}, 0},
	}
	for _, test := range tests {
		resp := UpdateNickNameService(test.req)
		if resp.Ret != test.ret {
			t.Errorf("UpdateNickNameService didn't pass. username:%s, ret:%d", test.req.UserName, test.ret)
		}
	}
}
