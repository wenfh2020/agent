package main

import (
	"agent/common"
	"fmt"
)

// func Base64Encode(src []byte) []byte {
// 	return []byte(base64.StdEncoding.EncodeToString(src))
// }

// func Base64Decode(src []byte) []byte {
// 	return []byte(base64.StdEncoding.DecodeString(string(src)))
// }

func test_md5() {
	s := "1234"
	fmt.Printf("%X\n", common.Md5String(s))
}

func test_base64() {
	fmt.Println("test_base64...")

	data := "hello world12345!?$*&()'-@~"
	fmt.Printf("src str: %v\n", data)

	enstr := common.Base64Encode(data)
	fmt.Printf("base64 encode str: %v\n", enstr)

	// Base64 Standard Decoding
	destr, err := common.Base64Decode(enstr)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}

	fmt.Printf("base64 decode str: %v\n", string(destr))
}

func main() {
	test_base64()
	// test_md5()
}
