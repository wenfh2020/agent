package main

import (
	"agent/common"
	"agent/db"
	"agent/inet"
	"agent/proto"
	"net/http"
	"runtime"

	"github.com/spf13/viper"
	log "github.com/thinkboy/log4go"
)

func initHTTP() (err error) {
	log.Trace("InitHTTP %v", viper.GetString("net.http.addrs"))

	httpServerMux := http.NewServeMux()
	httpServerMux.HandleFunc("/product/agent/check", proto.AgentCheck)

	log.Info("start http listen:\"%s\"", viper.GetString("net.http.addrs"))

	network, addr, err := inet.ParseNetwork(viper.GetString("net.http.addrs"))
	if err != nil {
		log.Error("ParseNetwork() error(%v)", err)
		return
	}

	go inet.HttpListen(httpServerMux, network, addr)
	return
}

func initDb() error {
	db.InitDbMgr()
	config := viper.Sub("mysql.mysql_lhl_product")
	return db.AddDbInfo("mysql_lhl_product", config.GetString("url"),
		config.GetInt("max_idle_conn"), config.GetInt("max_open_conn"))
}

func initConfig() error {
	viper.SetConfigFile("./conf/config.yml")
	return viper.ReadInConfig()
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

	// test_insert_db()
	// test_update_db()
	// test_select_db()

	common.InitSignal()
}
