package types

type Students struct {
	Id    int64 `json:"id"`
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"gte=1,lte=120"`
	Email string `json:"email" validate:"required,email"`

	/*
	validate:"required" → field must be present and non-empty
	validate:"gte=1,lte=120" → Age must be between 1 and 120
	validate:"email" → must be a valid email format
	*/
}
