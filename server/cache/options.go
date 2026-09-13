package cache

func WithBlockOverwrite() CacheOption {
	return func(c *Cache) {
		c.blockOverwrite = true
	}
}
