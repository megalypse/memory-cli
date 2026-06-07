package keyterm

func Extract(raw string) []string {
	return GetRegex().FindAllString(raw, -1)
}
