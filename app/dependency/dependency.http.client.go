package dependency

import (
	"resty.dev/v3"
)

func ProvideHttpClient() *resty.Client {
	var (
		r = resty.New()
	)
	return r
}
