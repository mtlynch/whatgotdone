package validate

func Email(email string) bool {
	if len(email) <= 2 {
		return false
	}
	return true
}