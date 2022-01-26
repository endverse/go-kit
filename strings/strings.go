package strings

import (
	"crypto/rand"
	"io"
	"strings"
)

// idLen is a length of captcha id string.
// (20 bytes of 62-letter alphabet give ~119 bits.)
const idLen = 20

// idChars are characters allowed in captcha id.
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// rngKey is a secret key used to deterministically derive seeds for
// PRNGs used in image and audio. Generated once during initialization.
var rngKey [32]byte

func init() {
	if _, err := io.ReadFull(rand.Reader, rngKey[:]); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}

}

type randomString struct {
	str string
}

func (r *randomString) String() string {
	return r.str
}

func (r *randomString) ToLower() *randomString {
	return &randomString{
		str: strings.ToLower(r.str),
	}
}

func (r *randomString) ToUpper() *randomString {
	return &randomString{
		str: strings.ToUpper(r.str),
	}
}

// RandomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-9. The slice can be used as a captcha
// solution.
func RandomDigits(length int) []byte {
	return randomBytesMod(length, 10)
}

// randomId returns a new random id string.
func RandomString(l int) *randomString {
	var length int
	if l > 0 {
		length = l
	} else {
		length = idLen
	}
	b := randomBytesMod(length, byte(len(idChars)))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return &randomString{
		str: string(b),
	}
}

func CutString(s string, i int) string {
	s = strings.ToLower(s)

	if strings.HasPrefix(s, "gray-") {
		i = 11
	}

	if len(s) <= i {
		return s
	}

	return s[:i]
}

func FirstLetterToLower(s string) string {
	if s == "" {
		return ""
	}

	if len(s) == 1 {
		return strings.ToLower(s)
	}

	return strings.ToLower(s[:1]) + s[1:]
}
