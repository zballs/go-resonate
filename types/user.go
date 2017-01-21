package types

const (
	// User types
	ARTIST   = "artist"
	LISTENER = "listener"
)

// User
type User struct {
	Email    string `json:"email"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	// Other info?
}

func NewUser(email, name, password, _type string) *User {
	return &User{
		Email:    email,
		Name:     name,
		Password: password,
		Type:     _type,
	}
}

func NewArtist(email, name, password string) *User {
	return &User{
		Email:    email,
		Name:     name,
		Password: password,
		Type:     ARTIST,
	}
}

func NewListener(email, name, password string) *User {
	return &User{
		Email:    email,
		Name:     name,
		Password: password,
		Type:     LISTENER,
	}
}
