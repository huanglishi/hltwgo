
## 配置文件
```
文件路径 /conf/app.conf 修改数据路径账号密码

```
## 开发、环境
  下载安装go环境>=1.16 （推荐最新1.19）
  
  mysql(建议版本8，字符集：utf8mb4/utf8mb4_general_ci)
## 开发新增路径访问权限配置
### 1请求是提示“403”
配置文件：在conf\app.conf文件的allowurl字段添加你的域名
### 2请求提示 “您的请求不合法，请按规范请求数据!”
#### 解决方式
（1）需要验证则在请求头添加 {"verify-time":当前时间戳,"verify-encrypt":md5(特定字符串+当前时间戳)},特定字符串和conf\app.conf文件下的apisecret值相同

  (2) 不需要验证的在配置文件：routers\route.go文件中的validityAPi()方法下 strings.Contains(c.Request.URL.Path, "您新增的访问路径")
### 3请求提示 “和token相关的错误!”
#### 解决方式
（1）如果您接口需要token验证则在请求头添加{Authorization:your token}

  (2) 不需要token验证（例如：登录）则在配置文件：utils\tool\jwt.go文件中的noVerify下添加“您新增的访问路径”
## 本地运行项目
 安装mysql数据库，导入数据文件（更目录下有hltw.sql，对应服务器上tuwen_saa_go）
## 安装fresh 热更新-边开发边编译
go install github.com/pilu/fresh@latest

## 初始化mod
go mod tidy

# 热编译运行
bee run 或 fresh 
# 打包
go build main.go
# 打包（此时会打包成Linux上可运行的二进制文件，不带后缀名的文件） 
#### 在项目根目录cmd进入执行
```
SET GOOS=linux
SET GOARCH=amd64
go build

```
在更目录生成运行的二进制文件：huling,拷贝文件到（49.234.109.66服务器->go项目->tuwensaas)