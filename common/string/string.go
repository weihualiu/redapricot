package string

func IsExistUpper(s string) bool {
	f := false
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x41 && s[i] <= 0x5A {
			f = true
			break
		}
	}
	return f
}
