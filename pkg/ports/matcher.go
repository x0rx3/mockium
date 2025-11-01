package ports

type MatcherMap[T Request] map[Matcher[T]]ResponseBuilder[T]

func NewMatcherMap[T Request]() MatcherMap[T] {
	return make(MatcherMap[T])
}

type Matcher[T Request] interface {
	Match(req T) bool
}
