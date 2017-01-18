package types

const (
	// User types
	ARTIST   = "artist"
	LISTENER = "listener"
)

// User
type User struct {
	Email string `json:"email"`
	Name  string `json:"username"`
	Type  string `json:"type"`
	// Other info?
}

func NewUser(email, name, _type string) *User {
	return &User{
		Email: email,
		Name:  name,
		Type:  _type,
	}
}

func NewArtist(email, name string) *User {
	return &User{
		Email: email,
		Name:  name,
		Type:  ARTIST,
	}
}

func NewListener(email, name string) *User {
	return &User{
		Email: email,
		Name:  name,
		Type:  LISTENER,
	}
}
