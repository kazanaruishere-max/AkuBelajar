package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with custom rules.
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator with custom validation rules.
func New() *Validator {
	v := validator.New()

	// Custom: Indonesian phone number (E.164 format)
	v.RegisterValidation("phone_id", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		matched, _ := regexp.MatchString(`^\+62\d{9,12}$`, phone)
		return matched
	})

	// Custom: NISN (10 digit numeric)
	v.RegisterValidation("nisn", func(fl validator.FieldLevel) bool {
		nisn := fl.Field().String()
		matched, _ := regexp.MatchString(`^\d{10}$`, nisn)
		return matched
	})

	// Custom: NIP (18 digit numeric)
	v.RegisterValidation("nip", func(fl validator.FieldLevel) bool {
		nip := fl.Field().String()
		matched, _ := regexp.MatchString(`^\d{18}$`, nip)
		return nip == "" || matched // optional field
	})

	// Custom: password strength
	v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		pw := fl.Field().String()
		if len(pw) < 8 || len(pw) > 72 {
			return false
		}
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
		hasDigit := regexp.MustCompile(`\d`).MatchString(pw)
		hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(pw)
		return hasUpper && hasLower && hasDigit && hasSpecial
	})

	return &Validator{validate: v}
}

// ValidateStruct validates a struct and returns formatted error messages.
func (v *Validator) ValidateStruct(s interface{}) []ValidationError {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   toSnakeCase(e.Field()),
			Message: formatErrorMessage(e),
		})
	}
	return errors
}

// ValidationError represents a single field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", toSnakeCase(e.Field()))
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", toSnakeCase(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", toSnakeCase(e.Field()), e.Param())
	case "strong_password":
		return "Password harus min 8 karakter, mengandung huruf besar, kecil, angka, dan simbol"
	case "phone_id":
		return "Nomor HP harus format +62xxxxxxxxx"
	case "nisn":
		return "NISN harus 10 digit angka"
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", toSnakeCase(e.Field()), e.Param())
	default:
		return fmt.Sprintf("%s tidak valid", toSnakeCase(e.Field()))
	}
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
