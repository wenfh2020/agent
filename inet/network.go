package inet

import (
	"agent/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"

	"net"
	"net/http"

	log "github.com/thinkboy/log4go"

	"time"
)

const (
	networkSpliter = "@"
)

type AckNormal struct {
	Errno  int    `json:"errno"`
	Errstr string `json:"errstr"`
}

func HttpListen(mux *http.ServeMux, network, addr string) {
	log.Trace("httpListen")

	httpServer := &http.Server{
		Handler:      mux,
		ReadTimeout:  viper.GetDuration("net.http.read_time_out"),
		WriteTimeout: viper.GetDuration("net.http.write_time_out"),
	}
	httpServer.SetKeepAlivesEnabled(true)
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Error("net.Listen(\"%s\", \"%s\") error(%v)", network, addr, err)
		panic(err)
	}

	if err := httpServer.Serve(l); err != nil {
		log.Error("server.Server() error(%v)", err)
		panic(err)
	}
}

func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, networkSpliter); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	} else {
		network = str[:idx]
		addr = str[idx+1:]
		return
	}
}

func ParseReqBody(req interface{}, w http.ResponseWriter, r *http.Request) (body string, errno int) {
	log.Trace("ParseReqBody")

	var err error
	var bytes []byte

	if bytes, err = ioutil.ReadAll(r.Body); err != nil {
		errno = common.ERR_INVALID_BODY
		log.Error("ioutil.ReadAll() failed (%s)", err)
		return
	}

	if err = json.Unmarshal(bytes, req); err != nil {
		errno = common.ERR_INVALID_BODY
		log.Error(err)
		return
	}

	body = string(bytes)
	log.Debug("msg body = %s", body)
	return
}

func sendErrorAck(w http.ResponseWriter, r *http.Request, errno *int, body *string, start time.Time) {
	log.Trace("sendErrorAck")

	ack := AckNormal{*errno, common.GetCodeMsg(*errno)}
	bytes, err := json.Marshal(ack)
	if err != nil {
		log.Error("invalid ack body")
		return
	}

	if _, err = w.Write(bytes); err != nil {
		log.Error("w.Write(\"%s\") error(%v)", bytes, err)
	}

	log.Info("req: \"%s\", post: \"%s\", res:\"%s\", ip:\"%s\", time:\"%fs\"", r.URL.String(), *body, bytes, r.RemoteAddr, time.Now().Sub(start).Seconds())
}

func SendAck(w http.ResponseWriter, r *http.Request, ack interface{}, errno *int, body *string, start time.Time) {
	log.Trace("SendAck")

	if *errno != 0 {
		sendErrorAck(w, r, errno, body, start)
		return
	}

	bytes, err := json.Marshal(ack)
	if err != nil {
		log.Error("invalid ack body")
		return
	}

	if _, err = w.Write(bytes); err != nil {
		log.Error("w.Write(\"%s\") error(%v)", bytes, err)
	}

	log.Info("req: \"%s\", post: \"%s\", res:\"%s\", ip:\"%s\", time:\"%fs\"", r.URL.String(), *body, bytes, r.RemoteAddr, time.Now().Sub(start).Seconds())
}
