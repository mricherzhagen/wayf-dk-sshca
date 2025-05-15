package main

import (
	_ "embed"
	"sshca"
)

var (
	//go:embed ca.template
	tmpl string
)

func main() {
	publicKey, signer := sshca.GetSignerFromSshAgent()

	sshca.Config = sshca.Conf{
		Template:                  tmpl,
		HTMLTemplate:              "login",
		Verification_uri_template: "http://localhost:2280/%s\n",
		SSOTTL:                    "3m",
		RendevouzTTL:              "1m",
		SshListenOn:               "localhost:2221",
		WebListenOn:               "localhost:2280",
		Cryptokilib:               "/usr/lib/pkcs11/libsofthsm2.so",
		Slot:                      "1241298034",
		CaConfigs: map[string]sshca.CaConfig{
			"softCa": {
				Name:      "Soft CA",
				Signer:    signer,
				PublicKey: publicKey,
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
