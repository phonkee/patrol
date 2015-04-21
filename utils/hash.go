package utils

import "code.google.com/p/go.crypto/bcrypt"

func HashPassword(password, secret string, cost int) (string, error) {
	sum := []byte(secret + password)
	body, err := bcrypt.GenerateFromPassword(sum, cost)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func VerifyHashedPassword(hash, password, secret string) (result bool, err error) {
	hb := []byte(hash)
	sum := []byte(secret + password)
	r := bcrypt.CompareHashAndPassword(hb, sum)
	if r != nil {
		if r == bcrypt.ErrMismatchedHashAndPassword {
			result, err = false, nil
		} else {
			result, err = false, r
		}
	} else {
		result, err = true, nil
	}
	return
}
