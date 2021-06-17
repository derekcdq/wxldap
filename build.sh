#!/bin/sh

echo "开始检查需要的依赖包...."

for i in `grep -E "github.com|golang.org" wxopld/*.go|awk '{print $2 }'|sed s/\"//g`
do
        echo "开始安装 $i"
        go get $i
done

for i in `grep "github.com" main.go|sed s/\"//g`
do
        echo "开始安装 $i"
        go get $i
done

echo "开始编译Linux版本.."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/wx2openldap_linux main.go && echo "二进制文件已经生成到了bin目录下"
echo "开始编译Windows版本.."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/wx2openldap_win32 main.go && echo "二进制文件已经生成到了bin目录下"
echo "开始编译Mac版本.."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/wx2openldap_mac main.go && echo "二进制文件已经生成到了bin目录下"
