package main

import (
	"agent/common"
	mysql "agent/db"
	"agent/inet"
	"agent/proto"
	"net/http"
	"runtime"

	"github.com/spf13/viper"
	log "github.com/thinkboy/log4go"
)

func initConfig() error {
	viper.SetConfigFile("./conf/config.yml")
	return viper.ReadInConfig()
}

func initDb() error {
	mysql.InitDbMgr()
	cf := viper.Sub("mysql.mysql_lhl_product")
	return mysql.AddDbInfo("mysql_lhl_product", cf.GetString("url"),
		cf.GetInt("max_idle_conn"), cf.GetInt("max_open_conn"))
}

func initHTTP() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/product/agent/check", proto.HttpReqAgentCheck)

	addrs := viper.GetString("net.http.addrs")
	network, addr, err := inet.ParseNetwork(addrs)
	if err != nil {
		log.Error("parse http addrs error! err: %v", err)
		return err
	}

	log.Info("start http listen: \"%s\"", addrs)
	go inet.HttpListen(mux, network, addr)
	return nil
}

func main() {
	if err := initConfig(); err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(viper.GetInt("base.maxproc"))
	log.LoadConfiguration(viper.GetString("base.log"))
	defer log.Close()

	if err := initHTTP(); err != nil {
		panic(err)
	}

	if err := initDb(); err != nil {
		panic(err)
	}

	common.InitSignal()
}
