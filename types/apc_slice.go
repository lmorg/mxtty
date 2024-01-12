package types

type ApcSlice []string

func (as *ApcSlice) Value(i int) string {
	if len(*as) <= i {
		return ""
	}
	return (*as)[i]
}
