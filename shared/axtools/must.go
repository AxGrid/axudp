package axtools

func MustBytes(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}
