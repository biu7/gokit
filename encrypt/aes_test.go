package encrypt

import "testing"

func TestAES(t *testing.T) {
	key := []byte("0123456789abcdef")
	data := []byte("hello world")
	encrypted, err := AesCBCEncryptBase64(data, key)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := AesCBCDecryptBase64(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(decrypted) != string(data) {
		t.Fatal("decrypted data not equal to original data")
	}
}
