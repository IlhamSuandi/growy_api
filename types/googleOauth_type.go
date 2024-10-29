package types

type GoogleUser struct {
	Sub            string
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Email          string
	Email_verified bool
}
