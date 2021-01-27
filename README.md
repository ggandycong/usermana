# 用户管理系统

## 功能需求

实现一个用户管理系统，用户可以登录、拉取和编辑他们的profiles。 

用户可以通过在Web页面输入username和password登录，backend系统负责校验用户身份。

成功登录后，页面需要展示用户的相关信息；否则页面展示相关错误。

成功登录后，用户可以编辑以下内容：

1. 上传profile picture

2. 修改nickname（需要支持Unicode字符集，utf-8编码）

用户信息包括：

1. username（不可更改）

2. nickname

3. profile picture

需要提前将初始用户数据插入数据库用于测试。确保测试数据库中包含10,000,000条用户账号信息。

## 设计要求

* 分别实现HTTP server和TCP server，主要的功能逻辑放在TCP server实现
* Backend鉴权逻辑需要在TCP server实现
* 用户账号信息必须存储在MySQL数据库。通过MySQL Go client连接数据库
* 使用基于Auth/Session Token的鉴权机制，Token存储在redis，避免使用JWT等加密的形式。
* TCP server需要提供RPC API，RPC机制希望自己设计实现
* Web server不允许直连MySQL、Redis。所有HTTP请求只处理API和用户输入，具体的功能逻辑和数据库操作，需要通过RPC请求TCP server完成
* 尽可能使用Go标准库
* 安全性
* 鲁棒性
* 性能

## 开发环境

* 操作系统：macOS Catalina 10.15.6
* Go:1.15.6
* Mysql: 8.0.22
* Redis: 6.0.10

## 设计简介

* 此项目主要包含三个部分

  * 用户登陆界面
  * http server的服务
  * tcp server的服务

  Note:主要是http erver和tcp server模块，其中最重要的是tcp server。

* **HTTP Server**负责处理HTTP请求，对用户数据进行简单处理并转发至TCP服务。

* **TCP Server**处理HTTP服务转发的请求并访问MySQL和redis。

* http server和tcp server 之间的通信主要是通过rpc协议。

* **rpc**通信协议实现。

  * rpc客户端：通过tcp 连接rpc服务器，发送不同的数据信息，从而调用不同的rpc服务。
  * rpc服务端：首先注册一些服务(本质服务对应的服务处理函数), 监听请求。 当收到rpc客户端数据时，解析信息调用对应的服务并且将处理完成的结果返回给客户端。

  

## 实现流程

**整体流程图**

![diagram1](./resource/diagram1.jpg)

**具体登陆获取信息请求**

![diagram2](./resource/diagram2.jpg)

## API接口

### 1.注册接口信息

| URL                          | 方法 |
| ---------------------------- | ---- |
| http://localhost:1088/signUp | POST |

**输入参数**

| 参数名   | 描述   | 可选 |
| -------- | ------ | ---- |
| username | 用户名 | 否   |
| password | 密码   | 否   |
| nickname | 昵称   | 是   |



### 2.登录接口信息

| URL                         | 方法 |
| --------------------------- | ---- |
| http://localhost:1088/login | POST |

**输入参数**

| 参数名   | 描述   | 可选 |
| -------- | ------ | ---- |
| username | 用户名 | 否   |
| password | 密码   | 否   |

### 3.获取用户信息接口信息

>需要在登录接口之后调用

| URL                           | 方法 |
| ----------------------------- | ---- |
| http://localhost:1088/profile | GET  |

**输入参数**

| 参数名   | 描述   | 可选 |
| -------- | ------ | ---- |
| username | 用户名 | 否   |

### 4.更改用户昵称接口信息

> 需要在登录接口之后调用

| URL                                  | 方法 |
| ------------------------------------ | ---- |
| http://localhost:1088/updateNickName | POST |

**输入参数**

| 参数名   | 描述   | 可选 |
| -------- | ------ | ---- |
| username | 用户名 | 否   |
| nickname | 新昵称 | 否   |

### 5.更改用户头像接口信息

> 需要在登录接口之后调用

| URL                                    | 方法 |
| -------------------------------------- | ---- |
| http://localhost:1088/uploadProfilePic | POST |

**输入参数**

| 参数名   | 描述         | 可选 |
| -------- | ------------ | ---- |
| username | 用户名       | 否   |
| image    | 头像图片路径 | 否   |

## 数据储存

### mysql设计

主要维护两张表，一张保存用户信息， 一张保存用户登陆信息。

更详细可以参考**mysql/usermane.sql**文件创建表格sql语句

#### 用户信息表

| Field     | Field        | Null | Key  | Default | Extra          |
| --------- | ------------ | ---- | ---- | ------- | -------------- |
| id        | bigint       | NO   | PRI  | NULL    | auto_increment |
| user_name | varchar(255) | NO   | UNI  |         |                |
| nick_name | varchar(255) | NO   |      |         |                |
| pic_name  | varchar(255) | YES  |      |         |                |

note: **pic_name**字段主要用于存储用户头像的路径。

#### 用户登陆信息表

| Field     | Field        | Null | Key  | Default | Extra          |
| --------- | ------------ | ---- | ---- | ------- | -------------- |
| id        | bigint       | NO   | PRI  | NULL    | auto_increment |
| user_name | varchar(255) | NO   | UNI  |         |                |
| password  | varchar(255) | NO   |      |         |                |

### redis设计

redis缓存数据设计。

主要是缓冲登陆校验的token和用户信息，其中用户信息键值对中的值，是一个哈希表，表中有三项元素，分表是代表用户信息是否有效，用的的nick_name, 用户的pic_name。

| key           | value                                          |
| ------------- | ---------------------------------------------- |
| auto_username | Token                                          |
| username      | { [valid, 1/""],[nick_name, “”] [pic_name,“”]} |

## 代码结构

```bash
usermana
├── README.md								//程序文档
├── benchmark								//压力测试文件
├── config									//配置文件
├── httpServer							//http server
├── log											//日志相关文件	
├── mysql										//mysql
├── protocol								//主要定义一些通讯的数据结构
├── redis										//redis相关文件
├── resource								//文档所需要资源
├── rpc											//rpc实现
├── static									//用户头像存放路径
├── tcpServer								//tcp server
├── templates								//用户UI相关html
└── utils										//相关辅助函数

```

## 部署

1. 在conf/bench_conf配置系统资源
2. 在redis/config.go配置redis和mysql
3. 运行TCP server

```bash
cd tcpServer
go vet
go build tcpServer.go
./tcpServer
```

4. 运行HTTP server

```bash
cd httpServer
go vet
go build httpServer.go
./httpServer
```



## 功能测试

**用户登陆**

![login](./resource/login.jpg)

**用户登陆成功**

![login_success](./resource/login_success.jpg)

**显示用户信息**

> 默认用户图像信息为空

![updatePicName](./resource/updatePicName.jpg)

**修改用户头像路径**

![updatePicName_succes](./resource/updatePicName_succes.jpg)

## 单元测试

简单对redis,mysql,tcpserver接口进行单元测试.

分别在**redis, mysql, tcpServer**，目录下执行

```bash
go test
```

主要执行**redis/redis_test.go**，**mysql/mysql_test.go**,**tcpServer/tcpServer_test.go**测试文件.

## 压力测试

### 模拟

使用程序模拟多用户并发请求，一个协程(goroutine)对应一个用户，并发启动多个协程，即可达到模拟测试目的。

运行**benchmark/benchmark.go**文件，即可模拟压力测试。

### login:

#### 固定用户 200 并发

> Total Requests(50000) - Concurrency(200) - Random(false) - Cost(2.785033353s) - QPS(17954/sec)

#### 随机用户 200 并发

> Total Requests(50000) - Concurrency(200) - Random(true) - Cost(2.934093639s) - QPS(17042/sec)

#### 固定用户 2000 并发

> Total Requests(50000) - Concurrency(2000) - Random(false) - Cost(3.485402234s) - QPS(14346/sec)

#### 随机用户 2000 并发

> Total Requests(50000) - Concurrency(2000) - Random(true) - Cost(4.493010268s) - QPS(11129/sec)

### profile:

#### 固定用户 200 并发

> Total Requests(5000) - Concurrency(200) - Random(false) - Cost(133.75811ms) - QPS(37381/sec)

#### 随机用户 200 并发

> Total Requests(5000) - Concurrency(200) - Random(true) - Cost(130.668901ms) - QPS(38265/sec)

#### 固定用户 2000 并发

> Total Requests(5000) - Concurrency(2000) - Random(false) - Cost(383.837833ms) - QPS(13027/sec)

#### 随机用户 2000 并发

> Total Requests(5000) - Concurrency(2000) - Random(true) - Cost(404.677492ms) - QPS(12356/sec)

### updateNickName:

#### 固定用户 200 并发

> Total Requests(5000) - Concurrency(200) - Random(false) - Cost(140.719206ms) - QPS(35532/sec)

#### 随机用户 200 并发

> Total Requests(5000) - Concurrency(200) - Random(true) - Cost(151.213109ms) - QPS(33066/sec)

#### 固定用户 2000 并发

> Total Requests(5000) - Concurrency(2000) - Random(false) - Cost(376.417359ms) - QPS(13284/sec)

#### 随机用户 2000 并发

> Total Requests(5000) - Concurrency(2000) - Random(true) - Cost(424.903766ms) - QPS(11768/sec)

## 问题总结

这边整理一下遇到的问题，以及需要注意的点。

1. 使用mysql/mysql_test.go去初始化数据库，插入1 000 000条数据，这个过程比较久，我在本机测试，需要90分钟。而**go test**默认最多运行10分钟，当运行大于10分钟，程序自动退出。因此为了顺利完整插入数据，可以设置**timeout**参数，指定go test运行完timeout时间后才退出。

   ```bash
   go test -timeout=120m
   ```

   表示**go test**运行时，直到完成所有测试或者到达120分钟， 程序才会退出。

2. 在高并发多连接的情况下，还需要设置相关参数，比如一个进程最多允许打开的连接数，文件描述符等，还有mysql数据库连接池相关参数的设置。

   * **kern.ipc.somaxconn**：表示内核最多并发连接数。默认为128，推荐在1024-4096之间，数字越大占用内存也越大。

   ```bash
   sudo sysctl -w kern.ipc.somaxconn=2048
   ```

   * **kern.maxfilesperproc**：每个进程能够同时打开的最大文件数量(根据需要灵活配置)

   ```bash
   sudo sysctl -w kern.maxfiles=65536 
   ```

   * **kern.maxfiles**:  系统中允许的最多文件数量 （根据需要灵活配置）

   ```bash
   sudo sysctl -w kern.maxfiles=22288
   ```

   * **ulimit -n**:  打开的文件描述符(本系统默认2560), 这个设置受到maxfilesperproc和maxfiles的约束，不能大于其中任何一个。

   ```bash
   ulimit -n 10000
   ```

   Note: 每个控制台都需要独立设置才生效。

   * **max_connections**：数据库最多连接数， 默认是150。

   通过以下sql语句查看max_connections值。 

   ```sql
   show variables like "max_connections";
   ```

   通过一下sql语句可以设置（根据需要灵活配置）

   ```sql
   set global max_connections = 2000;
   ```

   当遇到连接数超过max_connections时候，会报类似如下错误：

   ```
   err:"Error 1040: Too many connections"
   ```

   * mysql数据库初始化一些连接配置

   ```go
   	//设置mysql每个连接的生命周期，默认值0，即连接不关闭一直存在。
     //此值的设置应该小于wait_timeout，避免连接因为超时而意外关闭。wait_timeout的默认值28800，8个小时，
   	//ConnMaxLifetime值推荐为wait_timeout的一半，即14400，4个小时。
   	db.SetConnMaxLifetime(config.ConnMaxLifetime)
     //MaxIdleConns 连接池中最大空闲连接数.
   	db.SetMaxIdleConns(config.MaxIdleConns)
     // MaxOpenConns 同时连接数据库中最多连接数.  MaxIdleConns和MaxOpenConns最好设置一样
   	db.SetMaxOpenConns(config.MaxOpenConns)
   ```

3. 重点是**back_log**参数

   此参数影响mysql server中tcp accept队列中的连接数。 当tcp accetp队列容量超过**back_log**时候，那么mysql server将会丢弃此连接，从而导致tcp 协议栈回应一个rst包给客户端。(back_log默认值151)

   客户端收到rst之后，表现类似如下:

   ```
   read tcp 127.0.0.1:55718->127.0.0.1:3306: read: connection reset by peer
   ```

   然而，在**go-sql-driver/mysql**中，客户端收到了rst包，还会继续重试maxBadConnRetries(默认2)次连接。

   因此，客户端还是有很大几率连接成功并且完成请求，除非maxBadConnRetries次连接都收到了rst包。

   不过，这对于整体请求来说，响应就变慢了。换言之，QPS会大大减少。

   因此，我们可以修改tcp accetp队列容量，使之比mysql的连接池容量还要大，就可以避免上述问题，从而提高QPS。

   tcp acctpt队列容量除了受到back_log制约，还受到系统参数kern.ipc.somaxconn的制约。

   换言之，tcp acctpt队列容量 = min(back_log,  kern.ipc.somaxconn)。

   back_log参数具体介绍可参考：[back_log](https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_back_log)

   **back_log**变量通过如下方式查看

   ```
   show global variables like '%back_log%';
   ```

   由于back_log是静态变量，因此通过在my.cnf文件修改，如下

   ```
   [mysqld]
   back_log = 1000
   ```

   Note: 这个主要是排查花了1天多的时间，还是没有找到问题的本质所在，最后是靠**Zhang Weiwei大佬**解决的。

   具体可以参考：[MySQL出现connection reset by peer问题排查](https://confluence.shopee.io/pages/viewpage.action?pageId=375667105)

   排查问题的思路也是非常值得学习。

