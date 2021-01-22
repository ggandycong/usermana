package mysql

import (
	// "math/rand"
	"strconv"
	"testing"
)

/*
// TestCreateAccount100 初始化测试数据库,创建10 000 000 个数据.
func TestCreateAccount100(t *testing.T) {
	for i := 0; i < 10000000; i++ {
		userName := "bot" + strconv.Itoa(i)
		if err := CreateAccount(userName, "1234"); err != nil {
			t.Errorf("CreateAccount didn't pass. username:%s, err:%q", userName, err)
		}
		if err := CreateProfile(userName, "bot"); err != nil {
			t.Errorf("CreateProfile didn't pass. username:%s, err:%q", userName, err)
		}
		if i%100 == 0 {
			//t.Log("now is %d", i)
		}
	}
}
*/

// TestCreateAccount 测试函数CreateAccount创建用户.
func TestCreateAccount(t *testing.T) {
	var tests = []struct {
		userName string
		password string
	}{
		{"botTest", "1234"},
	}
	for _, test := range tests {
		if err := CreateAccount(test.userName, test.password); err != nil {
			t.Errorf("CreateAccount didn't pass. username:%s, password:%s, err:%q", test.userName, test.password, err)
		}
	}
}

// TestCheckAccountExist 测试CheckAccountExist函数.
func TestCheckAccountExist(t *testing.T) {
	var tests = []struct {
		userName string
		exist    bool
	}{
		{"bot1", true},
		{"noExist", false},
	}
	for _, test := range tests {
		if ok, err := CheckAccountExist(test.userName); err != nil || ok != test.exist {
			t.Errorf("CheckAccountExist didn't pass. userName:%s, exist:%t", test.userName, test.exist)
		}
	}
}

// TestLoginAuth 测试登陆LoginAuth函数.
func TestLoginAuth(t *testing.T) {
	var tests = []struct {
		userName, password string
		ok                 bool
	}{
		{"bot1", "1234", true},
		{"bot2", "123", false},
		{"noExist", "1234", false},
		{"", "", false},
		{"bot1", "", false},
		{"", "123", false},
	}
	for _, test := range tests {
		if ok, err := LoginAuth(test.userName, test.password); err != nil || ok != test.ok {
			t.Errorf("LoginAuth didn't pass. userName:%s, password:%s, ok:%t", test.userName, test.password, test.ok)
		}
	}
}

// TestCreateProfile 测试创建用户信息CreateProfile函数.
func TestCreateProfile(t *testing.T) {
	var tests = []struct {
		userName string
		nickName string
	}{
		{"bot439", "botAAB"},
	}
	for _, test := range tests {
		if err := CreateProfile(test.userName, test.nickName); err != nil {
			t.Errorf("CreateProfile didn't pass. username:%s, nickName:%s, err:%q", test.userName, test.nickName, err)
		}
	}
}

// TestGetProfile 测试获取用户信息函数GetProfile.
func TestGetProfile(t *testing.T) {
	var tests = []struct {
		userName string
		hasData  bool
	}{
		{"bot1", true},
		{"bot2", true},
		{"noExist", false},
		{"", false},
	}
	for _, test := range tests {
		if _, _, ok, err := GetProfile(test.userName); err != nil || ok != test.hasData {
			t.Errorf("GetProfile didn't pass. userName:%s, hasdata:%t", test.userName, test.hasData)
		}
	}
}

// TestCheckProfileExist 测试CheckProfileExist函数.
func TestCheckProfileExist(t *testing.T) {
	var tests = []struct {
		userName string
		exist    bool
	}{
		{"bot1", true},
		{"noExist", false},
	}
	for _, test := range tests {
		if ok, err := CheckProfileExist(test.userName); err != nil || ok != test.exist {
			t.Errorf("CheckProfileExist didn't pass. userName:%s, exist:%t", test.userName, test.exist)
		}
	}
}

// TestUpdateProfile 测试更新用户信息函数UpdateProfile函数.
func TestUpdateProfile(t *testing.T) {
	var tests = []struct {
		userName, nickName, picName string
	}{
		{"bot2", "soy1234", ""},
	}
	for _, test := range tests {
		if _, err := UpdateProfile(test.userName, test.nickName, test.picName); err != nil {
			t.Errorf("UpdateProfile didn't pass. userName:%s, nickName:%s, picName:%s", test.userName, test.nickName, test.picName)
		}
	}
}

//TestUpdateNikcName 测试修改nickName函数.
func TestUpdateNikcName(t *testing.T) {
	var tests = []struct {
		userName, nickName string
	}{
		{"bot2", "soy12345"},
	}
	for _, test := range tests {
		if _, err := UpdateNikcName(test.userName, test.nickName); err != nil {
			t.Errorf("UpdateNikcName didn't pass. userName:%s, nickName:%s", test.userName, test.nickName)
		}
	}
}

//TestUpdateProfilePic 测试更新用户头像路径函数UpdateProfilePic.
func TestUpdateProfilePic(t *testing.T) {
	var tests = []struct {
		userName, picName string
	}{
		{"bot2", "http://127.0.0.1:1188/static/default.jpeg"},
	}
	for _, test := range tests {
		if _, err := UpdateProfilePic(test.userName, test.picName); err != nil {
			t.Errorf("UpdateProfilePic didn't pass. userName:%s, picName:%s", test.userName, test.picName)
		}
	}
}

func BenchmarkUpdateNikcName(b *testing.B) {
	// b.ReportAllocs()
	var tests = []struct {
		userName, nickName string
	}{
		{"bot2", "soy12345"},
	}
	for _, test := range tests {
		for i := 0; i < b.N; i++ {
			if _, err := UpdateNikcName(test.userName, test.nickName); err != nil {
				b.Errorf("UpdateNikcName didn't pass. userName:%s, nickName:%s", test.userName, test.nickName)
			}
		}
	}
}

//BenchmarkLoginSame 基准测试 用户登陆函数(对于同一个用户函数).
func BenchmarkLoginSame(b *testing.B) {
	// b.ReportAllocs()
	var tests = []struct {
		userName, password string
	}{
		{"bot2", "123"},
	}
	for _, test := range tests {
		for i := 0; i < b.N; i++ {
			if _, err := LoginAuth(test.userName, test.password); err != nil {
				b.Errorf("LoginAuth didn't pass. userName:%s, password:%s", test.userName, test.password)
			}
		}
	}
}

//BenchmarkLoginRadom 基准测试 用户登陆函数(随机用户).
func BenchmarkLoginRadom(b *testing.B) {
	// b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// if _, err := LoginAuth("bot"+strconv.Itoa(rand.Intn(10000000)), "123"); err != nil {
		if _, err := LoginAuth("bot"+strconv.Itoa(b.N), "123"); err != nil {
			b.Errorf("LoginAuth didn't pass.")
		}
	}

}
