package goda

type Chain[T interface{ IsZero() bool }] struct {
	value T
	Error error
}

func (c Chain[T]) ok() bool {
	return c.Error == nil && !c.value.IsZero()
}

func (c Chain[T]) MustGet() T {
	if c.Error != nil {
		panic(c.Error)
	}
	return c.value
}

func (c Chain[T]) GetError() error {
	return c.Error
}

func (c Chain[T]) GetOrElse(other T) T {
	if c.Error != nil {
		return other
	}
	return c.value
}

func (c Chain[T]) GetOrElseGet(other func() T) T {
	if c.Error != nil {
		return other()
	}
	return c.value
}

func (c Chain[T]) GetResult() (T, error) {
	return c.value, c.Error
}

func (c Chain[T]) IsZero() bool {
	return c.Error == nil && c.value.IsZero()
}

func (c Chain[T]) getError(e *error) T {
	if *e == nil {
		*e = c.Error
	}
	return c.value
}
