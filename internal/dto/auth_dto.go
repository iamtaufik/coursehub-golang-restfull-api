package dto


type RegisterDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}