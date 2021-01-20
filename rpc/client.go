package rpc

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
)

//RPCClient rpc客户端 包含一个连接池(连接rpc 服务器).
type RPCClient struct {
	pool chan net.TCPConn
}

//Client 创建connections个tcp连接, 连接到address中，并且将连接保存到连接池作为返回值返回.
func Client(connections int, address string) (RPCClient, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return RPCClient{}, err
	}

	//创建connections个连接，并将其保存到pool连接池中.
	pool := make(chan net.TCPConn, connections)
	for i := 0; i < connections; i++ {
		//laddr 本地地址默认.
		conn, err := net.DialTCP("tcp4", nil, tcpAddr)
		if err != nil {
			return RPCClient{}, errors.New("rpc: init client failed")
		}
		pool <- *conn
	}
	return RPCClient{pool: pool}, nil
}

//Call 对外提供方法， 调用服务端方法， resp必须为指针类型，保存返回结果数据.
func (r *RPCClient) Call(name string, req interface{}, resp interface{}) error {
	return r.call(name, req, resp)
}

//call 真正rpc调用逻辑，  使用rpc调用函数name(req), 并将结果保存到resp中.
func (r *RPCClient) call(name string, req interface{}, resp interface{}) error {
	//从连接池获取一个空闲连接.
	conn := r.getConn()
	defer r.releaseConn(conn)

	//对请求进行封装.
	reqBytes, err := r.packRequest(name, req)
	if err != nil {
		return err
	}
	//将数据发送到rpc服务器
	conn.Write(reqBytes)

	//首先读取数据包的大小，保存到n.
	dataLen := make([]byte, PackMaxSize)
	n, err := conn.Read(dataLen)
	if err != nil && err != io.EOF {
		return err
	}
	if n <= 0 {
		return errors.New("rpc.Call: not data")
	}

	//将一个dataLen转为64位int，保存到len中.
	len, err := strconv.ParseInt(string(dataLen[:PackMaxSize]), 10, 64)
	if err != nil {
		return err
	}
	//创建长度为len的字符数组buff，准备接收应答数据.
	buff := make([]byte, len)
	//读取长度为len的数据到buff中
	n, err = conn.Read(buff)
	if err != nil && err != io.EOF {
		return err
	}
	if n <= 0 {
		return errors.New("rpc.Call: not data")
	}

	//解析json数据buff，保存到resp数据结构中.
	if err = r.unpackResponse(resp, buff); err != nil {
		return err
	}
	return nil

}

// getConn 从RPCClient连接池中随机取出一个空闲连接。 如果没有，那将会阻塞.
func (r *RPCClient) getConn() (conn net.TCPConn) {
	select {
	case conn := <-r.pool:
		return conn
	}
}

//releaseConn 将连接conn重新写入到连接池中.
func (r *RPCClient) releaseConn(conn net.TCPConn) {
	select {
	case r.pool <- conn:
		// TODO 返回空值怎么可行？
		return
	}
}

//packRequest 对请求数据进行json封装，然后返回其对应的字符数组.
/*
	以字符串拼接[name,v]
	封装之前
	eg:
		name: Login
		(v)protocol.ReqLogin{
			UserName: userName,
			Password: password,
		}

	封装之后，以字符串的形式返回如下内容。
	body的长度
	{
		Name:  函数名的字符串
		Data：
			{
				UserName: userName,
				Password: password,
			}
	}
*/
func (r *RPCClient) packRequest(name string, v interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	reqSt := request{Name: name, Data: dataBytes}
	reqBytes, err := pack(reqSt)
	if err != nil {
		return nil, err
	}
	return reqBytes, nil
}

//unpackResponse 将respBytes数据拆包，保存到resp中.
func (r *RPCClient) unpackResponse(resp interface{}, respBytes []byte) (err error) {
	if err := json.Unmarshal(respBytes, resp); err != nil {
		return err
	}
	return nil
}
