package assets

var assets map[string][]byte

func init() {
	assets = make(map[string][]byte)
}

func Get(name string) []byte {
	return assets[name]
}
