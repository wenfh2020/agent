package proto

import (
	"agent/common"
	mysql "agent/db"
	"agent/inet"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/thinkboy/log4go"
)

type ClientInfo struct {
	Level string `json:"level"`
	Type  string `json:"type"`
}

type DeviceInfo struct {
	Mac     string `json:"mac"`
	Version string `json:"version"`
}

type ReqAgentCheck struct {
	Client ClientInfo `json:"client"`
	Device DeviceInfo `json:"device"`
	Time   string     `json:"time"`
	Sign   string     `json:"sign"`
}

type AckAgentCheck struct {
	Errno      int        `json:"errno"`
	Errstr     string     `json:"errstr"`
	Device     DeviceInfo `json:"device"`
	Activation string     `json:"activation"`
	Time       string     `json:"time"`
	Sign       string     `json:"sign"`
}

var ErrNoDbRecord = errors.New("record not found")

const (
	aesKey = "*#^AQaabTuabMK*%"
	salt   = "fkja98374dsf%$^#DFGDS%@@@SDFdrgt"
)

func HttpReqAgentCheck(w http.ResponseWriter, r *http.Request) {
	log.Trace("proto: HttpReqAgentCheck")

	var (
		errno int
		body  string
		req   ReqAgentCheck
		ack   AckAgentCheck
	)

	ack.Errno = common.ERR_OK
	defer inet.SendAck(w, r, &ack, &errno, &body, time.Now())

	body, errno = inet.ParseReqBody(&req, w, r)
	if errno != common.ERR_OK {
		log.Error("parse req body failed!")
		return
	}

	/* check sign. */
	time1 := req.Time
	sign := req.Sign
	mac := req.Device.Mac
	sign2 := common.SignEncode(salt, mac, time1)
	log.Debug("\nsign:  %X\nsign2: %X\n", sign, sign2)

	if strings.Compare(sign, sign2) != 0 {
		log.Error("invalid sign! \nmac: %v, \nsign:  %X\nsign2: %X\n", mac, sign, sign2)
		errno = common.ERR_INVALID_SIGN
		return
	}

	/* base64 decode. */
	decode, _ := common.Base64Decode(mac)

	/* aes decrypt mac. */
	decrypt, _ := common.AesCBCDecrypt([]byte(decode), []byte(aesKey))
	mac = string(decrypt)

	log.Debug("decrypt mac str: %v", string(mac))

	db, err := mysql.GetDB("mysql_lhl_product")
	if err != nil {
		log.Error("get db conn failed! db name: %s, err: %v", "mysql_lhl_product", err)
		errno = common.ERR_INVALID_DATA
		return
	}

	var info mysql.DeviceInfo

	/* check mac, select database. */
	if err = db.Where("device_mac = ?", mac).Find(&info).Error; err != nil {
		log.Error("select db from failed! device: %v", mac)
		errno = common.ERR_DB_NO_RECORD
		return
	}

	/* check is vaild. */
	if info.Status == 0 {
		errno = common.ERR_INVALID_STATUS
		return
	}

	activation := info.Activation

	/* check activation */
	if len(activation) == 0 {
		/* generate activation, it's md5 data. */
		format := fmt.Sprintf("%s+%s", mac, salt)
		actmd5 := fmt.Sprintf("%x", common.Md5String(format))
		log.Info("activation: %v", actmd5)

		/* aes encrypt. */
		encrypt, _ := common.AesCBCEncrypt([]byte(actmd5), []byte(aesKey))
		/* base64 encode. */
		activation = common.Base64Encode(string(encrypt))

		/* update db activation. */
		tx := db.Begin()
		update := map[string]interface{}{"activation": actmd5}
		if err = tx.Model(&info).Where("device_mac = ?", mac).Updates(update).Error; err != nil {
			log.Error("update device info failed! err: %v", err)
			errno = common.ERR_DB_UPDATE_FAILED
			return
		}
		if err = tx.Commit().Error; err != nil {
			log.Error("update device info failed! err: %v", err)
			tx.Rollback()
			errno = common.ERR_DB_UPDATE_FAILED
			return
		}

		log.Info("update activation, mac: %s, activation: %s", mac, activation)
	} else {
		/* aes encrypt. */
		encrypt, _ := common.AesCBCEncrypt([]byte(activation), []byte(aesKey))
		/* base64 encode. */
		activation = common.Base64Encode(string(encrypt))
	}

	now := time.Now()
	time2 := now.Format("2006-01-02 15:04:05")
	encrypt, _ := common.AesCBCEncrypt([]byte(mac), []byte(aesKey))
	mac = common.Base64Encode(string(encrypt))
	sign = common.SignEncode(salt, mac, time2)

	ack.Activation = activation
	ack.Device.Version = req.Device.Version
	ack.Device.Mac = mac
	ack.Time = time2
	ack.Sign = sign
	return
}
