package service

type Comparer interface {
	Compare(expected, actual any) bool
}
