package cache

type Setter interface {
	Set(string, []byte) bool
}

type Getter interface {
	Get(string) ([]byte, bool)
}

type Deleter interface {
	Delete(string)
}

type Flusher interface {
	Flush()
}
