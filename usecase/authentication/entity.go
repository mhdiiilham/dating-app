package authentication

// SignUpRequest struct holds the data definition required for a user to sign up.
type SignUpRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

// SignUpResponse struct hold the data definition when signing up an user is success.
type SignUpResponse struct {
	ID          string
	Email       string
	AccessToken string
}

// Credential struct holds the data definition required for a user to sign-in.
type Credential struct {
	Email    string
	Password string
}
