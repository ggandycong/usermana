package main

import (
	"usermana/config"
	"usermana/log"
	"usermana/mysql"
	"usermana/protocol"
	"usermana/redis"
	"usermana/rpc"
	"usermana/utils"
)

func main() {
	//init log.
	if err := log.Config(config.TCPServerLogPath, log.LevelInfo); err != nil {
		panic(err)
	}
	//init server.
	server := rpc.Server()
	//注册服务.
	panicIfErr(server.Register("SignUp", SignUp, SignUpService))
	panicIfErr(server.Register("Login", Login, LoginService))
	panicIfErr(server.Register("GetProfile", GetProfile, GetProfileService))
	panicIfErr(server.Register("UpdateProfilePic", UpdateProfilePic, UpdateProfilePicService))
	panicIfErr(server.Register("UpdateNickName", UpdateNickName, UpdateNickNameService))

	//监听并且处理连接.
	server.ListenAndServe(config.TCPServerAddr)
}

// panicIfErr 错误包裹函数.
func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// SignUp 注册接口.
func SignUp(v interface{}) interface{} {
	return SignUpService(*v.(*protocol.ReqSignUp))
}

// Login 登录接口.
func Login(v interface{}) interface{} {
	return LoginService(*v.(*protocol.ReqLogin))
}

// GetProfile 获取信息接口.
func GetProfile(v interface{}) interface{} {
	return GetProfileService(*v.(*protocol.ReqGetProfile))
}

// UpdateProfilePic 更新头像接口.
func UpdateProfilePic(v interface{}) interface{} {
	return UpdateProfilePicService(*v.(*protocol.ReqUpdateProfilePic))
}

// UpdateNickName 更新昵称接口
func UpdateNickName(v interface{}) interface{} {
	return UpdateNickNameService(*v.(*protocol.ReqUpdateNickName))
}

// SignUpService 注册接口的实际服务，同时用于在注册时向rpc传递参数类型.
func SignUpService(req protocol.ReqSignUp) (resp protocol.RespSignUp) {
	if req.UserName == "" || req.Password == "" {
		resp.Ret = 1
		return
	}
	if req.NickName == "" {
		req.NickName = req.UserName
	}

	if err := mysql.CreateAccount(req.UserName, req.Password); err != nil {
		resp.Ret = 2
		log.Errorf("tcp.signUp: mysql.CreateAccount failed. usernam:%s, err:%q", req.UserName, err)
		return
	}
	if err := mysql.CreateProfile(req.UserName, req.NickName); err != nil {
		resp.Ret = 2
		log.Errorf("tcp.signUp: mysql.CreateProfile failed. usernam:%s, err:%q", req.UserName, err)
		return
	}

	resp.Ret = 0
	return
}

// LoginService 登录接口的实际服务，同时用于在注册时向rpc传递参数类型.
func LoginService(req protocol.ReqLogin) (resp protocol.RespLogin) {
	ok, err := mysql.LoginAuth(req.UserName, req.Password)
	if err != nil {
		resp.Ret = 2
		log.Errorf("tcp.login: mysql.LoginAuth failed. usernam:%s, err:%q", req.UserName, err)
		return
	}
	//账号或密码不正确.
	if !ok {
		resp.Ret = 1
		return
	}
	token := utils.GetToken(req.UserName)
	err = redis.SetToken(req.UserName, token, int64(config.TokenMaxExTime))
	if err != nil {
		resp.Ret = 2
		log.Errorf("tcp.login: redis.SetToken failed. usernam:%s, token:%s, err:%q", req.UserName, token, err)
		return
	}
	resp.Ret = 0
	resp.Token = token
	log.Infof("tcp.login: login done. username:%s", req.UserName)
	return
}

// GetProfileService 获取信息接口的实际服务，同时用于在注册时向rpc传递参数类型.
func GetProfileService(req protocol.ReqGetProfile) (resp protocol.RespGetProfile) {
	// 校验token
	ok, err := checkToken(req.UserName, req.Token)
	if err != nil {
		resp.Ret = 3
		log.Errorf("tcp.getProfile: checkToken failed. usernam:%s, token:%s, err:%q", req.UserName, req.Token, err)
		return
	}
	if !ok {
		resp.Ret = 1
		return
	}

	// 先尝试从redis取数据.
	nickName, picName, hasData, err := redis.GetProfile(req.UserName)
	if err != nil {
		resp.Ret = 3
		log.Errorf("tcp.getProfile: redis.GetProfile failed. username:%s, err:%q", req.UserName, err)
		return
	}
	if hasData {
		log.Infof("redis tcp.getProfile done. username:%s", req.UserName)
		return protocol.RespGetProfile{Ret: 0, UserName: req.UserName, NickName: nickName, PicName: picName}
	}

	//redis没有数据，从mysql里取.
	nickName, picName, hasData, err = mysql.GetProfile(req.UserName)
	if err != nil {
		resp.Ret = 3
		log.Errorf("mysql tcp.getProfile: mysql.GetProfile failed. username:%s, err:%q", req.UserName, err)
		return
	}
	if hasData {
		// 向redis插入数据.
		redis.SetNickNameAndPicName(req.UserName, nickName, picName)
	} else {
		resp.Ret = 2
		log.Errorf("tcp.getProfile: mysql.GetProfile can't find username. username:%s", req.UserName)
		return
	}
	log.Infof("tcp.getProfile done. username:%s", req.UserName)
	return protocol.RespGetProfile{Ret: 0, UserName: req.UserName, NickName: nickName, PicName: picName}

}

// UpdateProfilePicService 更新头像接口的实际服务(picName/FileName)，同时用于在注册时向rpc传递参数类型.
func UpdateProfilePicService(req protocol.ReqUpdateProfilePic) (resp protocol.RespUpdateProfilePic) {
	// 校验token.
	ok, err := checkToken(req.UserName, req.Token)
	if err != nil {
		resp.Ret = 3
		log.Errorf("tcp.updateProfilePic: checkToken failed. username:%s, token:%s, err:%q", req.UserName, req.Token, err)
		return
	}
	if !ok {
		resp.Ret = 1
		return
	}

	// 使redis对应的数据失效（由于数据将会被修改）.
	if err := redis.InvaildCache(req.UserName); err != nil {
		resp.Ret = 3
		log.Errorf("tcp.updateProfilePic: redis.InvaildCache failed. username:%s, err:%q", req.UserName, err)
		return
	}
	// 写入数据库.
	ok, err = mysql.UpdateProfilePic(req.UserName, req.FileName)
	if err != nil {
		resp.Ret = 3
		log.Errorf("tcp.updateProfilePic: mysql.UpdateProfilePic failed. username:%s, filename:%s, err:%q", req.UserName, req.FileName, err)
		return
	}
	if !ok {
		resp.Ret = 2
		return
	}
	resp.Ret = 0
	log.Infof("tcp.updateProfilePic done. username:%s, filename:%s", req.UserName, req.FileName)
	return
}

// UpdateNickNameService 更新昵称接口的实际服务(NickName)，同时用于在注册时向rpc传递参数类型.
func UpdateNickNameService(req protocol.ReqUpdateNickName) (resp protocol.RespUpdateNickName) {
	// 校验token.
	ok, err := checkToken(req.UserName, req.Token)
	if err != nil {
		resp.Ret = 3
		log.Errorf("tcp.updateNickName: checkToken failed. username:%s, token:%s, err:%q", req.UserName, req.Token, err)
		return
	}
	if !ok {
		resp.Ret = 1
		return
	}
	// 使redis对应的数据失效（由于数据将会被修改）.
	if err := redis.InvaildCache(req.UserName); err != nil {
		resp.Ret = 3
		log.Errorf("tcp.updateNickName: redis.InvaildCache failed. username:%s, err:%q", req.UserName, err)
		return
	}
	// 写入数据库.
	ok, err = mysql.UpdateNikcName(req.UserName, req.NickName)
	if err != nil {
		resp.Ret = 3
		log.Errorf("tcp.updateNickName: mysql.UpdateNikcName failed. username:%s, nickname:%s, err:%q", req.UserName, req.NickName, err)
		return
	}
	if !ok {
		resp.Ret = 2
		return
	}
	resp.Ret = 0
	log.Infof("tcp.updateNickName done. username:%s, nickname:%s", req.UserName, req.NickName)
	return
}

//checkToken  检查Token
func checkToken(userName string, token string) (bool, error) {
	// 压测token
	if token == "test" {
		return true, nil
	}
	return redis.CheckToken(userName, token)
}
