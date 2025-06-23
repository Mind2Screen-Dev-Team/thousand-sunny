package dependency

import (
	"encoding/base64"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xsecurity"
)

// # NOTE
//
//	Please Generate With Command:
//	`openssl rand -base64 32`
func ProvideEncryptionAES(c config.Cfg) (*xsecurity.EncryptionAES, error) {
	key, err := base64.StdEncoding.DecodeString(c.Security.AESKey["default"])
	if err != nil {
		return nil, err
	}

	enc, err := xsecurity.NewEncryptionAES(key)
	if err != nil {
		return nil, err
	}

	return enc, nil
}
