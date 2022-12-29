package action

import "reflect"

type Fn[T any] func(T)

type ActionSet[T any] struct {
	actions map[uintptr]Fn[T]
	removed []uintptr
	clear   bool
}

func NewSet[T any]() *ActionSet[T] {
	return &ActionSet[T]{
		actions: map[uintptr]Fn[T]{},
		removed: nil,
	}
}

func (actionSet *ActionSet[T]) Add(fn Fn[T]) Fn[T] {
	ptr := reflect.ValueOf(fn).Pointer()
	actionSet.actions[ptr] = fn
	return fn
}

func (actionSet *ActionSet[T]) Remove(fn Fn[T]) {
	ptr := reflect.ValueOf(fn).Pointer()
	actionSet.removed = append(actionSet.removed, ptr)
}

func (actionSet *ActionSet[T]) ClearNextApply() {
	actionSet.clear = true
}

func (actionSet *ActionSet[T]) Apply(x T) {
	for _, fn := range actionSet.actions {
		fn(x)
	}
	if actionSet.clear {
		for ptr := range actionSet.actions {
			delete(actionSet.actions, ptr)
		}
		actionSet.removed = actionSet.removed[:0]
		actionSet.clear = false
	} else {
		for _, ptr := range actionSet.removed {
			delete(actionSet.actions, ptr)
		}
		actionSet.removed = actionSet.removed[:0]
	}
}
