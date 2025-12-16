package goda

type Chain[T interface{ IsZero() bool }] struct {
	value  T
	eError error
	eType  string
	eFunc  string
}

func (c Chain[T]) ok() bool {
	return c.eError == nil && !c.value.IsZero()
}

func (c Chain[T]) MustGet() T {
	if c.eError != nil {
		panic(c.eError)
	}
	return c.value
}

func (c Chain[T]) GetError() error {
	return c.eError
}

func (c Chain[T]) GetOrElse(other T) T {
	if c.eError != nil {
		return other
	}
	return c.value
}

func (c Chain[T]) GetOrElseGet(other func() T) T {
	if c.eError != nil {
		return other()
	}
	return c.value
}

func (c Chain[T]) GetResult() (T, error) {
	return c.value, c.eError
}

func (c Chain[T]) IsZero() bool {
	return c.eError == nil && c.value.IsZero()
}

func (c Chain[T]) mergeError(e *error) T {
	if *e == nil {
		*e = c.eError
	}
	return c.value
}

func (c *Chain[T]) enterFunction(typeName string, funcName string) bool {
	if !c.ok() || c.eType != "" {
		return false
	}
	c.eType = typeName
	c.eFunc = funcName
	return true
}

func (c *Chain[T]) leaveFunction(flag bool) {
	if !flag {
		return
	}
	if c.eError == nil {
		c.eType = ""
		c.eFunc = ""
		return
	}
	//goland:noinspection GoTypeAssertionOnErrors
	if ce, ok := c.eError.(*Error); ok {
		ce.typeName = c.eType
		ce.funcName = c.eFunc
	}
}
