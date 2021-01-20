package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"usermana/log"
)

//serverFunc 处理实际请求的函数.
type serverFunc func(interface{}) interface{}

//request 对应rpc客户端请求的数据.
type request struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type rpcHandler struct {
	handler    serverFunc
	argsType   reflect.Type //handler函数的参数类型.
	replysType reflect.Type //handler函数的返回值类型.
}

// RPCServer 维护函数名以及函数具柄的map集合.
type RPCServer struct {
	router map[string]rpcHandler
}

//Server 初始化并返回一个rpc服务端.
func Server() RPCServer {
	return RPCServer{make(map[string]rpcHandler)}
}

//Register 注册服务端方法，服务端需实现两个函数，其中handler用于获取句柄，service用于获取实际参数类型.
func (r *RPCServer) Register(name string, handler serverFunc, service interface{}) error {
	return r.register(name, handler, service)
}

// ListenAndServe 服务端开始监听并响应请求.
func (r *RPCServer) ListenAndServe(address string) error {
	//监听.
	listener, err := r.listen(address)
	if err != nil {
		return err
	}
	//阻塞接收连接，并且处理该连接.
	err = r.accept(listener)
	if err != nil {
		return err
	}
	return nil
}

//register 注册服务，name函数名，函数对应的handler.
func (r *RPCServer) register(name string, handler serverFunc, service interface{}) error {
	//通过反射，获取service的类型。(eg:main.SignUpService).
	serviceType := reflect.TypeOf(service)
	//检查handleType的类型的各种属性.
	err := r.checkHandlerType(serviceType)
	if err != nil {
		return err
	}
	//获取handleType的参数以及返回值类型.
	argsType := serviceType.In(0)
	replysType := serviceType.Out(0)
	//将对应的[name,rpcHandler]保存起来.
	r.router[name] = rpcHandler{handler: handler, argsType: argsType, replysType: replysType}
	return nil
}

//checkHandlerType 检查handlerType类型的各种属性值.
func (r *RPCServer) checkHandlerType(handlerType reflect.Type) error {
	// 判断是否是函数类型.
	if handlerType.Kind() != reflect.Func {
		return errors.New("rpc.Register: handler is not func")
	}
	// 判断参数数量.
	if handlerType.NumIn() != 1 {
		return errors.New("rpc.Register: handler input parameters number is wrong, need one")
	}
	// 判断返回值数量.
	if handlerType.NumOut() != 1 {
		return errors.New("rpc.Register: handler output parameters number is wrong, need one")
	}
	// 判断参数和返回值类型.
	if handlerType.In(0).Kind() != reflect.Struct || handlerType.Out(0).Kind() != reflect.Struct {
		return errors.New("rpc.Register: parameters must be Struct")
	}
	return nil
}

//listen 监听address，并返回对应的具柄.
func (r *RPCServer) listen(address string) (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp4", laddr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

//accept 接受新的连接 启动协程RPCServer.handle去处理.
func (r *RPCServer) accept(listener *net.TCPListener) error {
	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		defer conn.Close()
		//启动协程 处理新的连接的业务逻辑
		go r.handle(conn)
	}
}

//handle 主要是读取rpc Client发过来的数据，并且将处理结果发送回去.
func (r *RPCServer) handle(conn *net.TCPConn) {
	if conn == nil {
		//log.Errorf("rpc.ListenAndServe: tcp connection is nil")
	}
	dataLen := make([]byte, PackMaxSize)
	for {
		//获取包的长度保存到n.
		n, err := conn.Read(dataLen)
		if err != nil && err != io.EOF {
			log.Errorf("rpc.ListenAndServer: connection read header failed. err:%q", err)
		}
		if n <= 0 {
			log.Errorf("rpc.ListenAndServe: not data")
		}
		len, err := strconv.ParseInt(string(dataLen[:PackMaxSize]), 10, 64)
		if err != nil {
			//log.Errorf("rpc.ListenAndServer: parseInt failed. err:%q", err)
		}
		//读取长度为len包的内容.
		buff := make([]byte, len)
		n, err = conn.Read(buff)
		if err != nil {
			//log.Errorf("rpc.ListenAndServer: connection read body failed. err:%q", err)
		}
		if n <= 0 {
			//log.Errorf("rpc.ListenAndServe: body not data")
		}

		//调度,处理实际的内容.
		rsp, err := r.dispatcher(buff)
		if err != nil {
			//log.Errorf("rpc.ListenAndServer: dispatch failed. err:%q", err)
		}
		//封装rsp的应答.
		rspBytes, err := r.packResponse(rsp)
		if err != nil {
			//log.Errorf("rpc.ListenAndServer: pack response failed. err:%q", err)
		}
		//将结果发送回去
		conn.Write(rspBytes)
	}
}

//dispatcher 查看req对应Name名字，并且找到name对应的handle处理.
func (r *RPCServer) dispatcher(req []byte) (interface{}, error) {
	// 解析接口名
	var cReq request
	if err := json.Unmarshal(req, &cReq); err != nil {
		return nil, err
	}
	//获取函数名对应的handle
	rh, ok := r.router[cReq.Name]
	if !ok {
		return nil, fmt.Errorf("rpc.ListenAndServe: can't find handler named %s", cReq.Name)
	}

	//解析参数类型， 根据此类型去接收data内容 保存到args.  args即使handle的实际参数.
	args := reflect.New(rh.argsType).Interface()
	if err := json.Unmarshal(cReq.Data, args); err != nil {
		return nil, err
	}
	// 由rpcHandler的具柄handler来处理对应的内容.
	return rh.handler(args), nil
}

func (r *RPCServer) packResponse(v interface{}) ([]byte, error) {
	return pack(v)
}
