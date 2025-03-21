package main

// convenient wrapper around *T to check against nil,
// choosing between two values based on existence, etc.
type Optional[T any] struct {
	val *T
}

// convenient wrapper around *T to check against nil,
// choosing between two values based on existence, etc.
//
// NOTE: use sparingly as it is not supported by the stdlib
// and you'd likely have to use *T everywhere else anyway.
//
// only use for the convenience methods where oneliners make
// more sense to use (e.g. assigning to struct fields)
func Option[T any](val *T) Optional[T] {
	return Optional[T]{
		val: val,
	}
}

// returns the value of the option (pointer) if it is not a nil pointer
//
// otherwise just returns the given value
func (o Optional[T]) UnwrapOr(val T) T {
	if o.val != nil {
		return *o.val
	}
	return val
}

// returns the option (pointer) if it is not a nil pointer
//
// otherwise just returns the given pointer as an Option
func (o Optional[T]) Or(val *T) Optional[T] {
	if o.val != nil {
		return o
	}
	return Option(val)
}
