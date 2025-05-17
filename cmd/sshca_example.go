package main

import (
	_ "embed"
	"os"
	"sshca"
)

var (
	//go:embed ca.template
	tmpl string
)

func main() {
	_, signer := sshca.GetSignerFromSshAgent()
	arg := os.Args[1]
	pubkey_ssh, err := os.ReadFile(arg)
	if err != nil {
		panic("Couldn't read pubkey")
	}

	sshca.Config = sshca.Conf{
		Template:                  tmpl,
		HTMLTemplate:              "login",
		Verification_uri_template: "http://localhost:2280/%s\n",
		SSOTTL:                    "3m",
		RendevouzTTL:              "1m",
		SshListenOn:               "localhost:2221",
		WebListenOn:               "localhost:2280",
		Cryptokilib:               "/usr/lib/pkcs11/libsofthsm2.so",
		Slot:                      "sshca_token", // This is not a "Slot" but a "label". See pkcs.go:46-54
		NoOfSessions:              1,
		CaConfigs: map[string]sshca.CaConfig{
			"softCa": {
				Name:      "Soft CA",
				Signer:    nil, // TODO: Replace with hsmSigner using findPrivatKey and known public key file
				PublicKey: string(pubkey_ssh),
				Settings: sshca.Settings{
					Ttl: 36 * 3600,
				},
			},
			"transport": {
				Signer: signer,
			},
		},
	}
	sshca.InitPKCS11("1234")

	sshca.Sshca()
}
