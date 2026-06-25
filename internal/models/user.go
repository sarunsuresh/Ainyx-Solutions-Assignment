package models

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	DOB   string `json:"dob"  validate:"required,datetime=2006-01-02"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"  validate:"required"`
	Email string `json:"email" validate:"required,email"`
	DOB   string `json:"dob"   validate:"required,datetime=2006-01-02"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"` 
	DOB   string `json:"dob"`
}

type UserWithAgeResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"` 
	DOB   string `json:"dob"`
	Age   int    `json:"age"`
}
