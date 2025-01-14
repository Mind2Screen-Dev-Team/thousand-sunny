package dependency

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xmail"
)

func ProvideXGmail(c config.Cfg) *xmail.XGmail {
	var (
		cfg  = c.SMTP["gmail"]
		xcfg = xmail.Config{
			SMTPHost:    cfg.Host,
			SMTPPort:    cfg.Port,
			FromName:    cfg.Credential.Name,
			FromAddress: cfg.Credential.Email,
			Username:    cfg.Credential.Username,
			Password:    cfg.Credential.Password,
		}
	)

	return xmail.NewXGmail(xcfg)
}
