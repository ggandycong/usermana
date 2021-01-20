package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

//Sha256 对密码passwd进行sha256编码, 然后将其转为字符串返回.
func Sha256(passwd string) string {
	rh := sha256.New()
	rh.Write([]byte(passwd))
	return hex.EncodeToString(rh.Sum(nil))
}

//GetToken 将src生成Token字符串，并将其返回.
func GetToken(src string) string {
	h := md5.New()
	h.Write([]byte(src + strconv.FormatInt(time.Now().Unix(), 10)))
	return hex.EncodeToString(h.Sum(nil))
}
