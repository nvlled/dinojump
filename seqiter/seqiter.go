package seqiter

type Iterator[T any] struct {
	Index int
	Items []T
}

func CreateSeqIterator[T any](items ...T) Iterator[T] {
	return Iterator[T]{
		Index: 0,
		Items: items,
	}
}

func (iter *Iterator[T]) Current() T {
	return iter.Items[iter.Index]
}

func (iter *Iterator[T]) Next() T {
	i := iter.Index
	iter.Index = (i + 1) % len(iter.Items)
	return iter.Items[iter.Index]
}

func (iter *Iterator[T]) Prev() T {
	i := iter.Index
	iter.Index = i - 1
	if iter.Index < 0 {
		iter.Index = len(iter.Items) - 1
	}
	return iter.Items[iter.Index]
}
