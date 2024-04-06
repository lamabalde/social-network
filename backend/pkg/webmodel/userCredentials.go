package webmodel

import (
	"net/mail"
)

// TODO: add unmarshalling (move them from datamethods.go)
type UserCredentials struct {
	UserName    string `json:"userName"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
	Gender      string `json:"gender,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Image       string `json:"image,omitempty"`
	AboutMe     string `json:"aboutMe,omitempty"`
}

func (u *UserCredentials) Validate() string {
	if IsEmpty(u.UserName) {
		return "username missing"
	}
	if IsEmpty(u.Email) {
		return "email missing"
	}

	// check email
	// mail.ParseAddress accepts also local domens e.g. witout .(dot)
	address, err := mail.ParseAddress(u.Email)
	if err != nil {
		return "wrong email"
	}
	u.Email = address.Address // in case of full address like "Barry Gibbs <bg@example.com>"
	// the regex allows only Internet emails, e.g. with dot-atom domain (https://www.rfc-editor.org/rfc/rfc5322.html#section-3.4)
	// if !regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`).Match([]byte(email)) {
	// 	return "wrong email"
	// }

	if IsEmpty(u.Password) {
		return "password missing"
	}
	if IsEmpty(u.DateOfBirth) {
		return "dateBirth missing"
	}
	if IsEmpty(u.Gender) {
		return "gender missing"
	}
	if IsEmpty(u.FirstName) {
		return "First name missing"
	}
	if IsEmpty(u.LastName) {
		return "Last name missing"
	}

	return ""
}

// type UserOnline struct {
// 	ID              string `json:"id"`
// 	UserName        string `json:"userName"`
// 	LastMessageDate string `json:"lastMessageDate"`
// }
