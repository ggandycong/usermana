package rpc

import (
	"encoding/json"
	"errors"
	"strconv"
)

// PackMaxSize tcp包header最大size.
const PackMaxSize int = 4

//pack 对类型v进行json封装，将封装后的json长度和json数据，拼接为字符数组并返回.
func pack(v interface{}) ([]byte, error) {
	//对v类型使用json封装
	rspBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	rspLen := len(rspBytes)
	rspLenStr := strconv.Itoa(rspLen)
	intLen := len(rspLenStr)

	if intLen > PackMaxSize {
		return nil, errors.New("rpc: package is out of size")
	}

	tb := make([]byte, PackMaxSize+rspLen)
	zerob := []byte("0")
	intLen--
	// header
	for i := PackMaxSize - 1; i >= 0; i-- {
		if intLen >= 0 {
			tb[i] = []byte(rspLenStr)[intLen]
			intLen--
		} else {
			tb[i] = zerob[0]
		}
	}
	//body
	for i := 0; i < rspLen; i++ {
		tb[PackMaxSize+i] = rspBytes[i]
	}

	return tb, nil
}
