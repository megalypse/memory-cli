package views

var lastPressedKey string

func SetLastPressedKey(key string) {
	lastPressedKey = key
}

func GetLastPressedKey() string {
	return lastPressedKey
}
