package main

import (
	sicher "github.com/dsa0x/sicher/cmd"
)

func main() {
	sicher.Execute()
}

// func mains() {
// 	_, cipher := sicher.Encrypt("6368616e676520746869732070617373776f726420746f206120736563726574", []byte("dev"))
// 	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")
// 	enc := []byte(cipher)
// 	res := sicher.Decrypt("6368616e676520746869732070617373776f726420746f206120736563726574", nonce, enc)
// 	fmt.Printf("Encrypted value: %x \nDecrypted value: %s\n", cipher, res)
// }
