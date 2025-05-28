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
		RelayingParty:             "", // I don't think this is used?
		Template:                  tmpl,
		HTMLTemplate:              "login",
		Verification_uri_template: "http://localhost:2280/%s\n",
		SSOTTL:                    "3m",
		RendevouzTTL:              "1m",
		SshListenOn:               "0.0.0.0:2221",
		WebListenOn:               "0.0.0.0:2280",
		Cryptokilib:               "/usr/lib/pkcs11/libsofthsm2.so",
		Slot:                      "sshca_token", // This is not a "Slot" but a "label". See pkcs.go:46-54
		NoOfSessions:              1,
		CaConfigs: map[string]sshca.CaConfig{
			"softCa": {
				Id:        "softCa",
				Name:      "Soft CA",
				Signer:    nil, // TODO: Replace with hsmSigner using findPrivatKey and known public key file
				PublicKey: string(pubkey_ssh),
				Settings: sshca.Settings{
					Ttl: 36 * 3600,
				},
				ClientID: "foo",
				Op: sshca.Opconfig{
					Userinfo:             "http://idp:5556/dex/userinfo",
					Device_authorization: "http://idp:5556/dex/auth",
					Token:                "http://idp:5556/dex/token",
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
