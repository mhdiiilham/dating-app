package authentication

// SignUpRequest struct holds the data definition required for a user to sign up.
type SignUpRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

// AccessResponse struct hold the data definition when signing up an user is success.
// TODO: need to refactor struct name!!!
type AccessResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

// Credential struct holds the data definition required for a user to sign-in.
type Credential struct {
	Email    string
	Password string
}
