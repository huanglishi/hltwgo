
## 配置文件
```
文件路径 /conf/app.conf 修改数据路径账号密码

```
## 开发环境
  下载安装go环境>=1.16 （推荐最新1.19）
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