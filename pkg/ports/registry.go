package ports

type Registry[T any] interface {
	Add(string, T)
	Delete(string) error
	Get(string) (T, bool)
	GetOrNil(string) T
	GetAll() map[string]T
}
