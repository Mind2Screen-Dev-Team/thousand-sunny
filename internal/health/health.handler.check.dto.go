package health

import "github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"

type (
	HealthResponseBody xresp.GeneralResponse[any, any]
)

type (
	HealthRequestInput   struct{}
	HealthResponseOutput struct {
		Body   HealthResponseBody
		Status int
	}
)
