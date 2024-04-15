package validate

import (
	"net/mail"

	"github.com/febriandani/backend-user-service/protogen/golang/users"
)

func ValidateUserRegistration(u *users.User) map[string]string {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return map[string]string{
			"en": "Incorrect email format",
			"id": "Format email salah",
		}
	}

	if u.Email == "" {
		return map[string]string{
			"en": "Email cannot be empty",
			"id": "Email tidak boleh kosong",
		}
	}

	if u.Username == "" {
		return map[string]string{
			"en": "Username cannot be empty",
			"id": "Username tidak boleh kosong",
		}
	}

	if u.Password == "" {
		return map[string]string{
			"en": "Password cannot be empty",
			"id": "Kata sandi tidak boleh kosong",
		}
	}

	if u.Repassword == "" {
		return map[string]string{
			"en": "Re-Password cannot be empty",
			"id": "Ulangi kata sandi tidak boleh kosong",
		}
	}
	return nil
}

func ValidateUserLogin(u *users.User) map[string]string {

	if u.Email == "" {
		return map[string]string{
			"en": "Email cannot be empty",
			"id": "Email tidak boleh kosong",
		}
	}

	if u.Password == "" {
		return map[string]string{
			"en": "Password cannot be empty",
			"id": "Kata sandi tidak boleh kosong",
		}
	}

	return nil
}
