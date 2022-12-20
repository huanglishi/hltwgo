
## 上传github
```
echo "# hltwgo" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/huanglishi/hltwgo.git
git push -u origin main
```
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