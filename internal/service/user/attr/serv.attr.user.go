package attr

type (
	UserCreate struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)
