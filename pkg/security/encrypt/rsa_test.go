package encrypt

import "testing"

func TestRSA(t *testing.T) {
	publicKey := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSQI+CDVUUkJ3iQlBJCCy9tXA+
/k2Jltb+QivvA4GHl2IKaV8V+Vx8+uaP5ddko2ifFmRv5ayrIp997U1RoxgUMgLb
P+wHvKJuHCj2wOaMJfD5WCNSfoj5Z7k5mL2JVTvym1ToFOFLh/TUjw9HELpB2KqE
iPKYPUHkAhzRWlwkUwIDAQAB
-----END PUBLIC KEY-----`
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCSQI+CDVUUkJ3iQlBJCCy9tXA+/k2Jltb+QivvA4GHl2IKaV8V
+Vx8+uaP5ddko2ifFmRv5ayrIp997U1RoxgUMgLbP+wHvKJuHCj2wOaMJfD5WCNS
foj5Z7k5mL2JVTvym1ToFOFLh/TUjw9HELpB2KqEiPKYPUHkAhzRWlwkUwIDAQAB
AoGACuLci8BB5zArwGEe242ZnvYrCIiI7WCzJUveb06u5pCn6I4W8mphFKB5M8uj
VvqSDNcJQaPL6khY4Au42waTzr3yUtHm+SiqoRGGw69y1o0iU2y437IFFK6DE0Js
YIYR8edyfQByDGPousYx6f/pXBVbvWgFx+TvJQMdSEvr+OECQQCb+NtUhwtAcIj0
20Y0q2tLfB4tIksRxZEA83g+/4t7ij85mJ255PGDb/lRmtoAaYXgNcQxiYYfESli
Ps3f+slRAkEA8AvpYkoyttW5vDmRncikOqqQjT0Fd743tl+6cyEYM9kv0aQK+bru
IoUSi7S2Wor8StXb8nFV0cpQpYT6a4kqYwJADQBekBo9tt5qnDaXEJbld8Jd5ZdB
OLcGUlT5Rg8ZOnAYux1csH1EIJe56bKFz5j8PudcdmCnXHmawITJtoz6MQJBAJt7
+7XQvbyz/1SuLWX4WngtlugFSI9hRJD6vLhqDVU4LsMp8HoF3M27SGH5scxADl8i
2w7U3nO/CjHMSrBw7mUCQEC7VSRq9SQuTKzFXR3tTi4jqnVcGXoaeFB2lDzTnwp+
EDS7qXxO8l5auhuFp6ylzWSFC0yvmCCoCJnbvqlNvg4=
-----END RSA PRIVATE KEY-----`

	src := "123"
	encryptData, err := RSAEncrypt([]byte(src), []byte(publicKey))
	if err != nil {
		t.Fatal(err)
	}

	dst, err := RSADecrypt(encryptData, []byte(privateKey))
	if err != nil {
		t.Fatal(err)
	}

	if src != string(dst) {
		t.FailNow()
	}
}
