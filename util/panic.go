package util

// PanicOnError will raise a panic if error is not nil
func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
