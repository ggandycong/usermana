package utils

import (
	"EntryTask/utils"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"path"
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

// GetFileName 为上传的文件生成一个文件名.
func GetFileName(fileName string, ext string) string {
	h := md5.New()
	h.Write([]byte(fileName + strconv.FormatInt(time.Now().Unix(), 10)))
	return hex.EncodeToString(h.Sum(nil)) + ext
}

// CheckAndCreateFileName 检查文件后缀合法性.
func CheckAndCreateFileName(oldName string) (newName string, isLegal bool) {
	ext := path.Ext(oldName)
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		//随机生成一个文件名.
		newName = utils.GetFileName(oldName, ext)
		isLegal = true
	}
	return newName, isLegal
}
