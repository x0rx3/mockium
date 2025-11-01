package ports

type Comparer interface {
	Compare(expected, actual any) bool
}
