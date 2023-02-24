package reg

import "regexp"

// GetNamedMap 提取命名字典
func GetNamedMap(pattern, target string) (result map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	if pattern == "" || target == "" {
		return
	}
	reg := regexp.MustCompile(pattern)
	m, n := reg.FindStringSubmatch(target), reg.SubexpNames()
	if len(m) > 1 && len(m) == len(n) {
		result = make(map[string]string)
		for i, v := range m {
			if i > 0 {
				result[n[i]] = v
			}
		}
	}
	return
}
