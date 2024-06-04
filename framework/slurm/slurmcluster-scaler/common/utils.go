package common

func AddIfNotPresent(s []string, target string) (bool, []string) {
	for _, ss := range s {
		if ss == target {
			return false, s
		}
	}
	return true, append(s, target)
}
