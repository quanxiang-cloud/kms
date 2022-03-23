package sign

import "crypto/md5"

// MD5 md5
func MD5(src []byte) ([]byte, error) {
	return MD5WithSalt(src, nil)
}

// MD5WithSalt md5 with salt
func MD5WithSalt(src []byte, salt []byte) ([]byte, error) {
	hash := md5.New()
	_, err := hash.Write(src)
	if err != nil {
		return nil, err
	}

	return hash.Sum(salt), nil
}
