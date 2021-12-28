package coder

import (
	"testing"
)

func TestCoder(t *testing.T) {
	testCase := struct {
		content       string
		encodeContent string
		decodeContent string
	}{
		content: `echo hello world
sdfasdSKFSdsfkjhasldkjfhklsdjhfkljahsdfkhaksdjhflkjashdflkjhaslkdjfh
askdjfhlkasjdhfklsjadhfkljsadhfkljshdlkfjhaklsdjhfklasjhdfkljahsdlkfjhsadlkfjhlkasdf
sdjkfhklsdahfkljhasdklfjhaskjdhfklajsdhfkjhasdkljhfkjasdhgjadfhgkjhadshfdsfkljhdskfjhadsfhkjadsxnfklhjdsajfhpsdaf;kdshjfkjhdslkfjhdsfjhklsadjfkhjwefhksadhjfuihewfjlksdjfklshjdfiuhawkejfwehiufghdsfghjklsadjhfuihsadfiujhsadojfkjasdhfiuadhsfpjhas;df
dsfkdsahlkfjhdsalkjfhkdlsjh
sleep 10
`,
		encodeContent: "H4sIAAAAAAAC/zyRsY7cQAxDf2VxX5DU16dJea0bIRqZpgR7EcJwPj/QePcqGjMi9ej5GH9wPDCqjsd1/C1fdnmY/Ov3ry9XJGEqTway5ERk0SCPhOU8qKQJ3tqz6Qwsu2l+VJro7RJtKvVWeGUQ9h1sIvy9YN7Jbu0Yj2Zj3iQ2M2DyrA5R3muMan3ddCxNjpXmgbXPXYhuVoRrArgCSXP92yMLdBkDT7nFZ7rAznNNElewAcwZCV4jkN2IcW4YV7D6r3Rf0GM7YVcOxjWwnbG2fQVvP9oi66nZ9OALNrbTHIpnt/js3s3rMrwQrOaDeIlYdtUYz8fPH8v+8T8AAP//lgPQKc8BAAA=",
		decodeContent: `echo hello world
sdfasdSKFSdsfkjhasldkjfhklsdjhfkljahsdfkhaksdjhflkjashdflkjhaslkdjfh
askdjfhlkasjdhfklsjadhfkljsadhfkljshdlkfjhaklsdjhfklasjhdfkljahsdlkfjhsadlkfjhlkasdf
sdjkfhklsdahfkljhasdklfjhaskjdhfklajsdhfkjhasdkljhfkjasdhgjadfhgkjhadshfdsfkljhdskfjhadsfhkjadsxnfklhjdsajfhpsdaf;kdshjfkjhdslkfjhdsfjhklsadjfkhjwefhksadhjfuihewfjlksdjfklshjdfiuhawkejfwehiufghdsfghjklsadjhfuihsadfiujhsadojfkjasdhfiuadhsfpjhas;df
dsfkdsahlkfjhdsalkjfhkdlsjh
sleep 10
`,
	}

	encodeContent, err := Encode(&testCase.content)
	if err != nil {
		t.Error(err)
		return
	}

	if encodeContent != testCase.encodeContent {
		t.Errorf("encode failed \n encode %s \n except %s", encodeContent, testCase.encodeContent)
	}

	var decodeContent string
	if err := Decode(encodeContent, &decodeContent); err != nil {
		t.Errorf("decode error %v", err)
	}

	if decodeContent != testCase.decodeContent {
		t.Errorf("decode failed \n decode %s \n except %s", decodeContent, testCase.decodeContent)
	}
}
