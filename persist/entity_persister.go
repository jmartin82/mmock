package persist

type EntityPersister interface {
	Read(name string) (string, error)
	Write(name, content string) error
	Delete(name string) error
	GetName() string
}
