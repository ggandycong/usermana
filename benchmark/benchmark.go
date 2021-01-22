package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func benchmarkBasicN(serverAddr string, n, c int32, isRan bool, ishttpPostMethod bool) (elapsed time.Duration) {
	readyGo := make(chan bool)
	//使用sync.WaitGroup等待线程结束.
	var wg sync.WaitGroup
	//表示需要等待的用户数量.
	wg.Add(int(c))

	remaining := n
	//为http请求创建的一个对象，用来保存多个请求过程中的一些状态.
	var transport http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          int(c),
		MaxIdleConnsPerHost:   int(c),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}
	//每个用户使用一个协程模拟请求
	cliRoutine := func(no int32) {
		//在函数退出时调用Done来通知wg，表示这个gorontinue已经完成.
		defer wg.Done()
		for atomic.AddInt32(&remaining, -1) >= 0 {
			// continue
			data := url.Values{}

			var buffer bytes.Buffer
			buffer.WriteString("bot")
			// rand
			if isRan {
				buffer.WriteString(strconv.Itoa(rand.Intn(10000000)))
			} else {
				buffer.WriteString("1")
			}

			username := buffer.String()

			data.Set("username", username)
			data.Set("password", "1234")
			data.Set("nickname", "newbot")
			//fmt.Printf("data = %s\n", data.Encode())
			var req *http.Request
			var err error
			if ishttpPostMethod {
				req, err = http.NewRequest("POST", serverAddr, bytes.NewBufferString(data.Encode()))
			} else {
				req, err = http.NewRequest("GET", serverAddr, bytes.NewBufferString(data.Encode()))
			}
			//设置http请求的cookie
			req.AddCookie(&http.Cookie{Name: "username", Value: username, Expires: time.Now().Add(120 * time.Second)})
			req.AddCookie(&http.Cookie{Name: "token", Value: "test", Expires: time.Now().Add(120 * time.Second)})

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
			if err != nil {
				log.Println(err)
			}
			<-readyGo
			//发送http 请求.
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
			}
			body, err1 := ioutil.ReadAll(resp.Body)
			if err1 != nil {
				log.Println(err1)
			}
			if len(body) <= 0 {
				log.Println("http response is empty.")
			}
			defer resp.Body.Close()
		}

	}
	//启动协程并发请求.
	for i := int32(0); i < c; i++ {
		go cliRoutine(i)
	}

	//关闭通道. (等待所有协程启动完，通过关闭管道通知协程发起http请求，这样可以让计时更加正确).
	close(readyGo)
	start := time.Now()
	//阻塞等待所有用户测试完毕.
	wg.Wait()

	return time.Since(start)
}

var num int64
var concurrency int64
var isRandom bool
var ishttpPostMethod bool

//init 初始化命令行参数默认值.
func init() {
	//用户数量.
	flag.Int64Var(&num, "n", 10000, "num")
	//并发数量.
	flag.Int64Var(&concurrency, "c", 2000, "concurrency")
	//测试用户是否随机(默认不随机).
	flag.BoolVar(&isRandom, "r", false, "isRandom")
	//请求方法(默认是get方法)，false对应的是get 方法.
	flag.BoolVar(&ishttpPostMethod, "p", false, "ishttpPostMethod")
}

func main() {
	//解析命令行参数.
	flag.Parse()
	//进行模拟测试.
	elapsed := benchmarkBasicN("http://127.0.0.1:1088/updateNickName", int32(num), int32(concurrency), isRandom, ishttpPostMethod)
	fmt.Println("HTTP server benchmark done:")
	fmt.Printf("\tTotal Requests(%v) - Concurrency(%v) - Random(%t) - Cost(%s) - QPS(%v/sec)\n",
		num, concurrency, isRandom, elapsed, math.Ceil(float64(num)/(float64(elapsed)/1000000000)))

}
