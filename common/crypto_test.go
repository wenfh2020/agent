package common

import (
	"agent/common"
	"fmt"
	"strings"
	"testing"
)

func TestMd5(t *testing.T) {
	src := "1234"
	dst := "81dc9bdb52d04dc20036dbd8313ed055"
	md5data := fmt.Sprintf("%X", common.Md5String(src))
	if strings.Compare(md5data, dst) != 0 {
		t.Errorf("invalid md5! md5data: %s, dst: %s", md5data, dst)
	}
}

func TestBase64(t *testing.T) {
	src := "hello world12345!?$*&()'-@~"
	encode := common.Base64Encode(src)
	dst, err := common.Base64Decode(encode)
	if err != nil {
		t.Errorf("base64 decode failed! err: %v", err)
	}
	if strings.Compare(src, dst) != 0 {
		t.Errorf("invalid base64! src: %s, dst: %s", src, dst)
	}
}
