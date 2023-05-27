package optional

func New[T any](t T) Optional[T] {
	return Optional[T]{val: t, present: true}
}

func NewEmpty[T any]() Optional[T] {
	return Optional[T]{present: false}
}

type Optional[T any] struct {
	val     T
	present bool
}

func (o *Optional[T]) IsPresent() bool {
	return o.present
}

func (o *Optional[T]) Get() T {
	return o.val
}

func (o *Optional[T]) GetOrDefault(t T) T {
	if o.present {
		return o.val
	}
	return t
}

func (o *Optional[T]) Set(t T) {
	o.val = t
	o.present = true
}
