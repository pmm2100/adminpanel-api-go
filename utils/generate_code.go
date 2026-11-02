package utils

import (
	"fmt"
	"math/rand/v2"
)

func GenerateEmailAuthenticationCode() string {
	var code string
	for i := 0; i < 6; i++ {
		letter := 'A' + rand.IntN(25)
		pickedLetter := fmt.Sprintf("%c", letter)
		code += pickedLetter
	}
	return code

}
