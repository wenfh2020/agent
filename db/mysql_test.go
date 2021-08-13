package db

/* go test mysql_test.go
 * go test mysql_test.go -v -run TestDbUpdate
 */

import (
	"agent/db"
	"agent/proto"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func initDb() error {
	db.InitDbMgr()
	config := viper.Sub("mysql.mysql_lhl_product")
	return db.AddDbInfo("mysql_lhl_product", config.GetString("url"),
		config.GetInt("max_idle_conn"), config.GetInt("max_open_conn"))
}

func initConfig() error {
	viper.SetConfigFile("../conf/config.yml")
	return viper.ReadInConfig()
}

func getDB(name string) (*gorm.DB, error) {
	if err := initConfig(); err != nil {
		return nil, err
	}

	if err := initDb(); err != nil {
		return nil, err
	}

	db, err := db.GetDB("mysql_lhl_product")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestDbInsert(t *testing.T) {
	db, err := getDB("mysql_lhl_product")
	if err != nil || db == nil {
		t.Errorf("get db failed! err: %v", err)
	}

	info := &proto.DbDeviceInfo{
		DeviceMac:     "AX-XX-XX-XX-XX-XX",
		DeviceVersion: "fdausfhuwhrw",
		Activation:    "fksurhuiwydjsf",
		ClientType:    "oiw3urfkdsnj",
		ClientLevel:   "dfsuyfewds",
		Status:        1,
		ActiveTime:    time.Now(),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err = db.Create(&info).Error; err != nil {
		t.Errorf("insert db record failed! err: %v", err)
	}
}

func TestDbUpdate(t *testing.T) {
	db, err := getDB("mysql_lhl_product")
	if err != nil {
		t.Errorf("get db failed! err: %v", err)
	}

	var info proto.DbDeviceInfo
	device := "XX-XX-XX-XX-XX-XX"
	tx := db.Begin()

	update := map[string]interface{}{
		"status": 0,
	}

	err = tx.Model(&info).Where("device_mac = ?", device).Updates(update).Error
	if err != nil {
		tx.Rollback()
		t.Errorf("update db failed! err: %v", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		t.Errorf("update db commit failed! err: %v", err)
	}
}

func TestDbSelect(t *testing.T) {
	db, err := getDB("mysql_lhl_product")
	if err != nil {
		t.Errorf("get db failed! err: %v", err)
	}

	var info proto.DbDeviceInfo
	device := "XX-XX-XX-XX-XX-XX"

	if err = db.Where("device_mac = ?", device).Find(&info).Error; err != nil {
		t.Errorf("select db failed! err: %v", err)
	}
}
