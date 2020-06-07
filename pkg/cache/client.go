package cache

type Client interface {
	Get(key string) (string, error)
	Set(key, value string) error
}
