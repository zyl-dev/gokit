package gochan

var gochanUUID int

func defaultUUID() int {
	gochanUUID++
	return gochanUUID
}
