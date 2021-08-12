/* https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go */
package main

import (
	"agent/common"
	"agent/proto"
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	aesKey := "*#^AQaabTuabMK*%"
	url := "http://127.0.0.1:7172/product/agent/check"
	salt := "fkja98374dsf%$^#DFGDS%@@@SDFdrgt"
	mac := "XX-XX-XX-XX-XX-XX"
	time := "2010-03-13 10:00:11"

	fmt.Printf("mac str: %v\n", mac)

	/* aes encrypt. */
	macEncrypt, _ := common.AesCBCEncrypt([]byte(mac), []byte(aesKey))
	fmt.Printf("mac aes encrypt hex: %X\n", macEncrypt)

	/* base64 encode. */
	mac = common.Base64Encode(string(macEncrypt))
	fmt.Printf("base64 encode str: %s\n", mac)

	/* digital signature */
	sign := common.SignEncode(salt, mac, time)

	var reqData proto.ReqAgentCheck
	reqData.Client.Level = "level vip"
	reqData.Client.Type = "product x3"
	reqData.Device.Version = "version 2.0.1.4"
	reqData.Device.Mac = mac
	reqData.Time = time
	reqData.Sign = sign

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(reqData)

	/* send request. */
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	/* response. */
	fmt.Println("response Status:", res.Status)
	fmt.Println("response Headers:", res.Header)
	bytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println("response Body:", string(bytes))

	var ack proto.AckAgentCheck

	/* check response. */
	if err = json.Unmarshal(bytes, &ack); err != nil {
		fmt.Println(err)
	}

	/* check sign. */
	time = ack.Time
	sign = ack.Sign
	mac = ack.Device.Mac
	sign2 := common.SignEncode(salt, mac, time)
	fmt.Printf("\nsign:  %X\nsign2: %X\n", sign, sign2)

	if strings.Compare(sign, sign2) != 0 {
		fmt.Println("invalid sign")
	}

	/* activation. */
	decode, _ := common.Base64Decode(ack.Activation)
	decrypt, _ := common.AesCBCDecrypt([]byte(decode), []byte(aesKey))
	activation := string(decrypt)
	fmt.Printf("activation: %v\n", activation)
}
