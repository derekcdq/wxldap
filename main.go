package main

import (
	"github.com/wonderivan/logger"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"wxldap/model"
)

func init() {
	model.InitLogConfig()
	model.InitConfig()
	new(model.WechatToken).Init()
	model.InitLdap()
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Info(sig)
		done <- true
	}()

	if os.Getppid() != 1 {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			logger.Info(err)
		}
		os.Exit(0)
	}
	//检查LDAP连接状态
	go func() {
		for {
			time.Sleep(time.Duration(10) * time.Second)
			model.ReConnLdap()
		}
	}()
	//刷新企业微信接口Token
	go func() {
		for {
			time.Sleep(time.Duration(1200) * time.Second)
			new(model.WechatToken).Init()
		}
	}()
	//同步部门及部门下人员
	go func() {
		for {
			model.SyncAllDept()
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()
    //保证进程不退出
	for {
		time.Sleep(time.Second * 1000)
	}

}