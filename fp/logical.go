package fp

func Not[A any](predicate func(A) bool) func(A) bool {
	return func(a A) bool {
		return !predicate(a)
	}
}

func And[A any](second func(A) bool) func(func(A) bool) func(A) bool {
	return func(first func(A) bool) func(A) bool {
		return func(a A) bool {
			return first(a) && second(a)
		}
	}
}

func Or[A any](second func(A) bool) func(func(A) bool) func(A) bool {
	return func(first func(A) bool) func(A) bool {
		return func(a A) bool {
			return first(a) || second(a)
		}
	}
}
