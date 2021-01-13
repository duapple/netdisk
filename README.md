# Netdisk程序使用说明



本程序是一款网盘服务程序，是采用web模式进行开发的，后端采用了够Go语言，前端使用JS，HTML和CSS。IDE使用的是vscode。

实现了最基本的网盘功能，包括：用户注册，登陆，上传文件，下载文件，上传文件夹，删除文件/文件夹，创建文件夹，目录跳转等，还有一些待实现的功能，如文件分享等。

### 使用方法：

---

克隆工程文件到`$GOPATH/src/github.com/user/`下。

```shell
$ git clone https://github.com/duapple/netdisk.git
```

然后运行`make`。

```shell
$ make
```

可能因为一些库不存在会报错，根据报错信息安装相应的库文件后， 重新编译。

编译成功后，修改配置文件`config/config.json`和`static/js/global.js`中的IP地址和端口号，可改为`localhost`和`9090`，然后执行可执行文件。

```shell
$ ./netdisk
```

在浏览器内输入`localhost:9090/login`进行登陆，目前注册仅支持用户名和密码，在注册页面输入用户名和密码完成注册，再登陆即可。进入Netdisk主页，新用户主页是空的，可以新建文件夹或者上传文件。

![网盘主页](https://github.com/duapple/Code/blob/master/pic/%E7%BD%91%E7%9B%98%E4%B8%BB%E9%A1%B5.png?raw=true)



### 工程说明

---

整个网盘程序目录结构如下：

```sh
.
├── DiskRoot
│   ├── INeedaGirl
│   ├── duapple
│   ├── duapple1
│   ├── hejiang
│   ├── q
│   ├── river
│   ├── river2
│   ├── river66
│   ├── shenhuifang
│   └── test
├── Makefile
├── README.md
├── config
│   ├── config.go
│   ├── config.json
│   └── log_config.go
├── database
│   ├── database.go
│   ├── user_db.json
│   └── user_db_map.json
├── ddns
│   └── setDDNS.go
├── homepage
│   ├── dir_opt.go
│   └── upload.go
├── loginpage
│   └── login.go
├── logrus.log
├── netdisk
├── netdisk.go
├── readme.txt
├── session
│   ├── main.go
│   ├── memory.go
│   └── session.go
├── static
│   ├── css
│   ├── imgs
│   ├── index.html
│   ├── js
│   ├── login.html
│   ├── regist.html
│   └── test.html
└── tmp
```



#### DiskRoot

用于存放用户的网盘数据。

#### config

存放Netdisk的一些配置数据，可以在程序运行前进行修改。

#### database

数据库，用户存储每个用户的个人信息数据，包括用户名密码，绑定的邮箱或者手机号等。目前没有引入数据库，仅仅采用json文件进行存储用户名和密码。

#### ddns

动态域名解析服务程序，目前域名用的是阿里云的，因为阿里云不支持在路由器端进行配置动态域名解析服务，所以这里写一个服务器端的服务程序来完成动态域名解析服务。当需要将Netdisk服务部署到公网上时使用。

#### homepage

主页功能代码，包括目录访问，上传下载等，采用Go的http包来完成基本功能，数据格式采用Formdata。文件下载续传采用http默认支持，文件上传采用分片处理方式来支持续传。

#### loginpage

登陆注册功能代码，包括用户注册和登陆功能的后端实现。

#### session

用于对cookie和sessin的支持。

#### static

web前端的所有代码。

#### tmp

文件上传时的缓存目录。分片上传时，用于文件缓存。

#### netdisk.go

Netdisk服务程序主函数。