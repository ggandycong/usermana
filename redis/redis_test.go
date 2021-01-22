package redis

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

// TestSetNickNameAndPicName 测试SetNickNameAndPicName函数.
func TestSetNickNameAndPicName(t *testing.T) {
	var tests = []struct {
		userName, nickName, picName string
	}{
		{"bot1", "soy1234", ""},
	}
	for _, test := range tests {
		if err := SetNickNameAndPicName(test.userName, test.nickName, test.picName); err != nil {
			t.Errorf("SetNickNameAndPicName didn't pass. userName:%s, nickName:%s, picName:%s, err:%q", test.userName, test.nickName, test.picName, err)
		}
	}
}

// TestGetProfile 测试Profile函数.
func TestGetProfile(t *testing.T) {
	var tests = []struct {
		userName string
		hasData  bool
	}{
		{"bot1", true},
		{"noExist", false},
		{"", false},
	}
	for _, test := range tests {
		if _, _, ok, err := GetProfile(test.userName); err != nil || ok != test.hasData {
			fmt.Printf("ok = %t, err = %q\n", ok, err)
			t.Errorf("GetProfile didn't pass. userName:%s, hasdata:%t, err:%q", test.userName, test.hasData, err)
		}
	}
}

//TestInvaildCache 测试InvaildCache函数.
func TestInvaildCache(t *testing.T) {
	var tests = []struct {
		userName string
	}{
		{"bot2"},
	}
	for _, test := range tests {
		if err := InvaildCache(test.userName); err != nil {
			t.Errorf("InvaildCache didn't pass. userName:%s, err:%q", test.userName, err)
		}
	}
}

// TestSetToken 测试SetToken函数.
func TestSetToken(t *testing.T) {
	var tests = []struct {
		userName string
		token    string
		exp      int64
	}{
		{"bot2", "auth", 5},
	}
	for _, test := range tests {
		if err := SetToken(test.userName, test.token, test.exp); err != nil {
			t.Errorf("SetToken didn't pass. userName:%s, token:%s, exp:%d, err:%q", test.userName, test.token, test.exp, err)
		}
	}
}

//TestCheckToken 测试CheckToken函数.
func TestCheckToken(t *testing.T) {
	var tests = []struct {
		userName string
		token    string
		ok       bool
	}{
		{"bot2", "auth", true},
		{"bot2", "auth2", false},
	}
	for _, test := range tests {
		if ok, err := CheckToken(test.userName, test.token); err != nil || ok != test.ok {
			t.Errorf("CheckToken didn't pass. userName:%s, token:%s, ok:%t, err:%q", test.userName, test.token, test.ok, err)
		}
	}
}

//BenchmarkSetTokenSame 基准测试SetToken函数(相同的用户名).
func BenchmarkSetTokenSame(b *testing.B) {
	// b.ReportAllocs()
	var tests = []struct {
		userName string
		token    string
		exp      int64
	}{
		{"bot2", "auth", 5},
	}
	for _, test := range tests {
		for i := 0; i < b.N; i++ {
			if err := SetToken(test.userName, test.token, test.exp); err != nil {
				b.Errorf("SetToken didn't pass. userName:%s, token:%s, exp:%d, err:%q", test.userName, test.token, test.exp, err)
			}
		}
	}
}

//BenchmarkSetTokenRandom 基准测试SetTokenRandom函数(用户名随机).
func BenchmarkSetTokenRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := SetToken("bot"+strconv.Itoa(rand.Intn(10000000)), "auth", 5); err != nil {
			b.Errorf("SetToken didn't pass")
		}
	}
}
