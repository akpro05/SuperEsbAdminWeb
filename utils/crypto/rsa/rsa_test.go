package rsa

import (
	"testing"
	//"tk.com/encoding/base64"
)

func TestPGP(t *testing.T) {
	cipher := "oiD6n7xW7vQ="

	tmp, err := base64.Decode([]byte(cipher))
	if err != nil {
		t.Error(err)
		return
	}

	pr, err := ReadPrivateKey("/home/dell3/tekno/Keystore/PGP/dev/private_key.cer")
	if err != nil {
		t.Error(err)
		return
	}
	plain, err := Decrypt(tmp, pr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(plain))
}
