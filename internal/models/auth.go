package models

type SignupRequest struct {
    Name     string `json:"name"     validate:"required"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    DOB      string `json:"dob"      validate:"required,datetime=2006-01-02"`
}

type LoginRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

type AuthUser struct {
    ID    int32  `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}