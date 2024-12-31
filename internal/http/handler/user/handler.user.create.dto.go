package handler_user

type (
	UserCreateDTO struct {
		Payload UserCreatePayload `in:"body=json"`
	}

	UserCreatePayload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)
