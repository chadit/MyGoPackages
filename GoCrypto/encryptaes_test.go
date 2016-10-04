package gocrypto

import (
	"encoding/base64"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	testmsg := "this is a test"
	eKey := "123456789012345678901234"
	var (
		r   []byte
		dr  []byte
		err error
	)

	if r, err = EncryptAES([]byte(eKey), []byte(testmsg)); err != nil {
		t.Fatalf("encrypt : %s", err)
	}

	if dr, err = DecryptAES([]byte(eKey), []byte(r)); err != nil {
		t.Fatalf("decrypt : %s", err)
	}
	//	fmt.Println(dr)
	if string(dr) != testmsg {
		t.Fatalf("expected %s but received %s", testmsg, dr)
	}
}

func TestAESEncrypt_Decrypt(t *testing.T) {
	testmsg := "this is a test"
	eKey := "123456789012345678901234"

	var (
		r   []byte
		dr  []byte
		err error
	)

	r = []byte{83, 206, 221, 78, 225, 94, 110, 119, 78, 247, 92, 234, 66, 72, 138, 3, 41, 92, 74, 109, 114, 75, 70, 216, 148, 119, 59, 176, 159, 240, 48, 162, 153, 158, 217, 101}
	sr := base64.StdEncoding.EncodeToString(r)
	data, _ := base64.StdEncoding.DecodeString(sr)

	if dr, err = DecryptAES([]byte(eKey), []byte(data)); err != nil {
		t.Fatalf("decrypt : %s", err)
	}
	if string(dr) != testmsg {
		t.Fatalf("expected %s but received %s", testmsg, dr)
	}
}

func TestAESEncrypt_WithID(t *testing.T) {
	testmsg := "5683e54571c9d4097cf8eca4"
	//	testmsg := "5683e54571c9d4097cf8eca4"
	eKey := "123456789012345678901234"
	var (
		r   []byte
		dr  []byte
		err error
	)

	if r, err = EncryptAES([]byte(eKey), []byte(testmsg)); err != nil {
		t.Fatalf("encrypt : %s", err)
	}

	if dr, err = DecryptAES([]byte(eKey), []byte(r)); err != nil {
		t.Fatalf("decrypt : %s", err)
	}
	if string(dr) != testmsg {
		t.Fatalf("expected %s but received %s", testmsg, dr)
	}
}
