package xmail

type XGmail struct {
	XMail
}

func NewXGmail(c Config) *XGmail {
	return &XGmail{
		XMail: New(c),
	}
}
