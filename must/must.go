package must

func MustString(text string, err error) string {
	if err != nil {
		panic(err)
	}
	return text
}
