package main

import (
	"fmt"
	"usermana/rpc"
)

// ReqLogin 登录请求.
type ReqLogin struct {
	UserName string `json:"user_name"` // 用户名, 不为空
	Password string `json:"password"`  // 密码, 不为空
}

func main() {
	rpcClient, err := rpc.Client(200, ":3194")
	if err != nil {
		panic(err)
	}

	req := ReqLogin{
		UserName: "username",
		Password: "passwd",
	}
	resp := ReqLogin{}
	if err := rpcClient.Call("Login", req, &resp); err != nil {
		fmt.Println("http.Login: Call Login failed.")
		return
	}

	fmt.Println(" done.")
}
