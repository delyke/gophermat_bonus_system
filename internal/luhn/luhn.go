package luhn

import "errors"

var (
	ErrInvalidChar = errors.New("invalid character")
	ErrNoDigits    = errors.New("no digits")
)

// Validate проверяет номер по алгоритму Луна.
// Вход: string или []byte (через []byte(s)).
// Разрешены цифры, пробелы и дефисы.
func Validate(s []byte) (bool, error) {
	sum, digits, err := sumLuhn(s, true)
	if err != nil {
		return false, err
	}
	if digits < 2 {
		return false, ErrNoDigits
	}
	return sum%10 == 0, nil
}

func sumLuhn(s []byte, includeLast bool) (sum, digits int, err error) {
	// Считает кол-во цифр
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			digits++
		case c == ' ' || c == '-':
			continue
		default:
			return 0, 0, ErrInvalidChar
		}
	}

	if digits == 0 {
		return 0, 0, ErrNoDigits
	}

	// Логика удвоения
	doubled := !includeLast

	// Идет с права на лево
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]

		if c < '0' || c > '9' {
			continue
		}

		d := int(c - '0')

		if doubled {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
		doubled = !doubled
	}

	return sum, digits, nil
}
