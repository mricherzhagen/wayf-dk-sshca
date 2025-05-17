package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	arg := os.Args[1]
	pubkey_ssh, err := os.ReadFile(arg)
	if err != nil {
		panic("Couldn't read pubkey")
	}
	pubkey, _, _, _, _ := ssh.ParseAuthorizedKey(pubkey_ssh)
	privkeyLabel := base64.RawURLEncoding.EncodeToString(pubkey.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)[:])
	fmt.Println(privkeyLabel)

}
