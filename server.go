package main

import (
	"agent/common"
	"agent/db"
	"agent/inet"
	"agent/proto"
	"errors"
	"net/http"
	"runtime"
	"time"

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
		config.GetInt("max-idle-conn"), config.GetInt("max-open-conn"))
	return nil
}

func initConfig() error {
	viper.SetConfigFile("./conf/config.yml")
	return viper.ReadInConfig()
}

func test_insert_db() {
	var info proto.DbDeviceInfo

	info.DeviceMac = "XX-XX-XX-XX-XX-XX"
	info.DeviceVersion = "fdausfhuwhrw"
	info.Activation = "fksurhuiwydjsf"
	info.ClientType = "oiw3urfkdsnj"
	info.ClientLevel = "dfsuyfewds"
	info.Status = 1
	info.ActiveTime = time.Now()
	info.CreateTime = time.Now()
	info.UpdateTime = time.Now()

	db, err := db.GetDB("mysql_lhl_product")
	if err != nil {
		panic(err)
	}

	if err = db.Create(&info).Error; err != nil {
		panic(err)
	}

	/* insert. */
	/* select. */
	/* update. */
}

func test_update_db() error {
	db, err := db.GetDB("mysql_lhl_product")
	if err != nil {
		panic(err)
	}

	var info proto.DbDeviceInfo
	device := "32132iuyeyruiq"
	tx := db.Begin()

	update := map[string]interface{}{
		"status": 0,
	}

	err = tx.Model(&info).Where("device_mac = ?", device).Updates(update).Error
	if err != nil {
		tx.Rollback()
		return errors.New("update fail")
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.New("commit failed")
	}
	return nil
}

func test_select_db() error {
	db, err := db.GetDB("mysql_lhl_product")
	if err != nil {
		panic(err)
	}

	var info proto.DbDeviceInfo
	device := "32132iuyeyrusiq"

	if err = db.Where("device_mac = ?", device).Find(&info).Error; err != nil {
		return err
	}
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

	// test_insert_db()
	// test_update_db()
	// test_select_db()

	common.InitSignal()
}
