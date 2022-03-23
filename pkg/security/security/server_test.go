package security

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestCommandString(t *testing.T) {
	if err := codeT("base64 std encode", "base64 std decode"); err != nil {
		t.Fatal(err)
	}

	if err := codeT("hex encode", "hex decode"); err != nil {
		t.Fatal(err)
	}

	server := New()
	dst, err := server.LookupCommandString("md5", []byte("123"))
	if err != nil {
		t.Fatal(err)
	}

	sortData := map[string]interface{}{
		"age":    18,
		"name":   "alex",
		"gender": true,
	}
	dst, err = server.LookupCommandString("sort", sortData)
	if err != nil {
		t.Fatal(err)
	}

	dst, err = server.LookupCommandString("md5|base64 std encode", []byte("123"))
	if err != nil {
		t.Fatal(err)
	}

	dst, err = server.LookupCommandString("aes cbc pkcs5 encode 123456789abcdefg", []byte("123"))
	if err != nil {
		t.Fatal(err)
	}

	dst, err = server.LookupCommandString(`rsa encode -----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSQI+CDVUUkJ3iQlBJCCy9tXA+
/k2Jltb+QivvA4GHl2IKaV8V+Vx8+uaP5ddko2ifFmRv5ayrIp997U1RoxgUMgLb
P+wHvKJuHCj2wOaMJfD5WCNSfoj5Z7k5mL2JVTvym1ToFOFLh/TUjw9HELpB2KqE
iPKYPUHkAhzRWlwkUwIDAQAB
-----END PUBLIC KEY-----`, []byte("123"))
	if err != nil {
		t.Fatal(err)
	}

	_ = dst
}

func codeT(cmd, unCmd string) error {
	var data = []byte("123")
	server := New()
	base64Std, err := server.LookupCommandString(cmd, data)
	if err != nil {
		return err
	}
	base64Std, err = server.LookupCommandString(unCmd, base64Std)
	if err != nil {
		return err
	}
	if string(data) != string(base64Std) {
		return errors.New("输入数据不等于结果数据")
	}

	return nil
}

func TestCommandSlice(t *testing.T) {
	security := []Security{
		{
			CMD: "md5",
			Arg: []string{"123"},
		},
		{
			CMD: "base64",
			Arg: []string{"std", "encode"},
		},
	}
	server := New()
	dst, err := server.LookUpCommandSilce(security, []byte("123"))
	if err != nil {
		t.Fatal(err)
	}

	_ = dst
}

func TestCommandwithJSON(t *testing.T) {
	cmdSting := `[{"cmd":"md5"},{"cmd":"base64","arg":["std","encode"]}]`
	security := make([]Security, 0, 2)
	err := json.Unmarshal([]byte(cmdSting), &security)
	if err != nil {
		t.Fatal(err)
	}

	server := New()
	dst, err := server.LookUpCommandSilce(security, []byte("123"))
	if err != nil {
		t.Fatal(err)
	}
	_ = dst
}

func TestLock(t *testing.T) {
	var rLock int32 = 0
	var wLock int32 = 1
	var lock int32 = rLock
	atomic.CompareAndSwapInt32(&lock, rLock, wLock)

	go func() {
		time.Sleep(time.Second * 10)
		atomic.CompareAndSwapInt32(&lock, wLock, rLock)
	}()
	for !atomic.CompareAndSwapInt32(&lock, rLock, rLock) {
	}

}

func TestQingCloudSign(t *testing.T) {

	data := struct {
		Count            int      `security:"count"`
		Vxnets1          []string `security:"vxnets"`
		Zone             string   `security:"zone"`
		InstanceType     string   `security:"instance_type"`
		SignatureVersion int      `security:"signature_version"`
		SignatureMethod  string   `security:"signature_method"`
		InstanceName     string   `security:"instance_name"`
		ImageID          string   `security:"image_id"`
		LoginMode        string   `security:"login_mode"`
		LoginPasswd      string   `security:"login_passwd"`
		Version          int      `security:"version"`
		AccessKeyID      string   `security:"access_key_id"`
		Action           string   `security:"action"`
		TimeStamp        string   `security:"time_stamp"`
	}{
		Count:            1,
		Vxnets1:          []string{"vxnet-0"},
		Zone:             "pek1",
		InstanceType:     "small_b",
		SignatureVersion: 1,
		SignatureMethod:  "HmacSHA256",
		InstanceName:     "demo",
		ImageID:          "centos64x86a",
		LoginMode:        "passwd",
		LoginPasswd:      "QingCloud20130712",
		Version:          1,
		AccessKeyID:      "QYACCESSKEYIDEXAMPLE",
		Action:           "RunInstances",
		TimeStamp:        "2013-08-27T14:30:10Z",
	}

	server := New()

	expect := `32bseYy39DOlatuewpeuW5vpmW51sD1A/JdGynqSpP8=`

	cmd := "sort query gonic asc|append begin GET\n/iaas/\n|sha256 SECRETACCESSKEY|base64 std encode"
	dst, err := server.LookupCommandString(cmd, &data)

	if err != nil {
		t.Fatal(err)
	}
	if string(dst) != expect {
		t.Errorf(`expect "%s" got "%s"`, expect, string(dst))
	}

	jsonData := `{"action":"RunInstances","access_key_id":"QYACCESSKEYIDEXAMPLE","count":1,"image_id":"centos64x86a","instance_name":"demo","instance_type":"small_b","login_mode":"passwd","login_passwd":"QingCloud20130712","signature_method":"HmacSHA256","signature_version":1,"time_stamp":"2013-08-27T14:30:10Z","version":1,"vxnets":["vxnet-0"],"zone":"pek1"}`
	dst, err = server.LookupCommandString(cmd, &jsonData)

	if err != nil {
		t.Fatal(err)
	}
	if string(dst) != expect {
		t.Errorf(`expect "%s" got "%s"`, expect, string(dst))
	}
}

// 32bseYy39DOlatuewpeuW5vpmW51sD1A_JdGynqSpP8=
