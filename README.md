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

## 实现流程

.....暂时略

## 数据储存

### mysql设计

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

redis缓存数据设计

| key           | value |
| ------------- | ----- |
| auto_username | Token |
| username      |       |

## 代码结构

## 部署

## 压测

### login:

#### 固定用户 200 并发

> Total Requests(5000) - Concurrency(200) - Random(false) - Cost(142.563168ms) - QPS(35073/sec)

#### 随机用户 200 并发

> Total Requests(5000) - Concurrency(200) - Random(true) - Cost(136.893744ms) - QPS(36525/sec)

#### 固定用户 2000 并发

> Total Requests(5000) - Concurrency(2000) - Random(false) - Cost(391.947248ms) - QPS(12757/sec)

#### 随机用户 2000 并发

> Total Requests(5000) - Concurrency(2000) - Random(true) - Cost(420.174593ms) - QPS(11900/sec)

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

