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

type AddressRequest struct {
    Line1      string `json:"line1"       validate:"required"`
    Line2      string `json:"line2"`
    City       string `json:"city"        validate:"required"`
    State      string `json:"state"       validate:"required"`
    PostalCode string `json:"postal_code" validate:"required"`
    Country    string `json:"country"     validate:"required"`
}

type UpdateProfileRequest struct {
    Name    string         `json:"name" validate:"required"`
    DOB     string         `json:"dob"  validate:"required,datetime=2006-01-02"`
    Address AddressRequest `json:"address" validate:"required"`
}

type AddressResponse struct {
    Line1      string `json:"line1"`
    Line2      string `json:"line2"`
    City       string `json:"city"`
    State      string `json:"state"`
    PostalCode string `json:"postal_code"`
    Country    string `json:"country"`
}

type ProfileResponse struct {
    ID      int32           `json:"id"`
    Name    string          `json:"name"`
    Email   string          `json:"email"`
    DOB     string          `json:"dob"`
    Age     int             `json:"age"`
    Address AddressResponse `json:"address"`
}