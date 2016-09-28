package gocrypto

import (
	"fmt"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	testmsg := "this is a test"
	eKey := "123456789012345678901234"
	var (
		r   string
		dr  string
		err error
	)

	if r, err = EncryptAES(eKey, testmsg); err != nil {
		t.Fatalf("encrypt : %s", err)
	}
	fmt.Println(r)
	if dr, err = DecryptAES(eKey, r); err != nil {
		t.Fatalf("decrypt : %s", err)
	}
	fmt.Println(dr)
	if dr != testmsg {
		t.Fatalf("expected %s but received %s", testmsg, dr)
	}

}
