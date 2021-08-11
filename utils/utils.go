package utils

func Uniq(str []string) []string {
	n := 0
	m := make(map[string]struct{}, len(str))
	for _, s := range str {
		if _, ok := m[s]; !ok {
			str[n] = s
			n++
			m[s] = struct{}{}
		}
	}
	return str[:n]
}

func StripSlash(str *string) {
	if len(*str) > 0 && (*str)[len(*str)-1] == '/' {
		*str = (*str)[:len(*str)-1]
	}
}

func Contains(target string, list []string) bool {
	for _, val := range list {
		if target == val {
			return true
		}
	}
	return false
}
