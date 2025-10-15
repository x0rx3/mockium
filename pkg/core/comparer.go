package core

type Comparer interface {
	Compare(expected, actual any) bool
}
