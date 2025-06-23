package dependency

import (
	"github.com/go-resty/resty/v2"
)

func ProvideHttpClient() *resty.Client {
	var (
		r = resty.New()
	)
	return r
}
