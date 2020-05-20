package loader

type Loader interface {
	ToURLsMap() (map[string]string, error)
}
