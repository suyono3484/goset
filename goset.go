package goset

import (
	"errors"
	"fmt"
	"reflect"
)

type Set struct {
	m map[any]any
}

type StrictSet struct {
	m   map[any]any
	typ reflect.Type
}

var ErrMismatchType = errors.New("must have same type for all the members")

func NewSet(v ...any) *Set {
	return NewSetSlice(v)
}

func NewSetSlice(v []any) *Set {
	mapper := make(map[any]any)
	for _, item := range v {
		mapper[item] = nil
	}

	m := &Set{
		m: mapper,
	}

	return m
}

func (s *Set) Add(v any) {
	mapper := s.m
	mapper[v] = nil
	s.m = mapper
}

func (s *Set) Remove(v any) {
	mapper := s.m
	delete(mapper, v)
	s.m = mapper
}

func (s *Set) Has(v any) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set) Len() int {
	return len(s.m)
}

func NewStrictSetErr(v ...any) (*StrictSet, error) {
	return NewStrictSetSliceErr(v)
}

func NewStrictSetSliceErr(v []any) (set *StrictSet, err error) {
	defer func() {
		if r := recover(); r != nil {
			set = nil
			err = r.(error)
		}
	}()

	set = NewStrictSetSlice(v)
	return
}

func NewStrictSet(v ...any) *StrictSet {
	return NewStrictSetSlice(v)
}

func NewStrictSetSlice(v []any) *StrictSet {
	if len(v) > 0 {
		mapper := make(map[any]any)
		typ := reflect.TypeOf(v[0])

		for _, item := range v {
			if reflect.TypeOf(item) != typ {
				panic(ErrMismatchType)
			}

			mapper[item] = nil
		}

		return &StrictSet{
			m:   mapper,
			typ: typ,
		}
	}

	return &StrictSet{}
}

func (s *StrictSet) AddErr(v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	s.Add(v)
	return
}

func (s *StrictSet) Add(v any) {
	if s.typ == nil {
		s.typ = reflect.TypeOf(v)
		s.m = map[any]any{v: nil}
		return
	}

	if reflect.TypeOf(v) != s.typ {
		panic(ErrMismatchType)
	}

	mapper := s.m

	mapper[v] = nil

	s.m = mapper
}

func (s *StrictSet) Remove(v any) {
	if _, ok := s.m[v]; ok {
		mapper := s.m

		delete(mapper, v)

		s.m = mapper
	}
}

func (s *StrictSet) Has(v any) bool {
	_, ok := s.m[v]
	return ok
}

func (s *StrictSet) Len() int {
	return len(s.m)
}

func (s *StrictSet) Distinct(v any) (int, error) {
	val := reflect.ValueOf(v)
	if val.IsNil() {
		return 0, errors.New("the output parameter is nil")
	}

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
		if val.Kind() != reflect.Slice {
			return 0, fmt.Errorf("unsupported type pointer of %v", val.Type())
		}
		if s.typ != val.Type().Elem() {
			return 0, fmt.Errorf("invalid output parameter type %v", val.Type().Elem())
		}

		for item := range s.m {
			val.Set(reflect.Append(val, reflect.ValueOf(item)))
		}
		return val.Len(), nil
	} else if val.Kind() == reflect.Slice {
		if len(s.m) > val.Len() {
			return 0, errors.New("insufficient length of the output parameter")
		}
		if s.typ != val.Type().Elem() {
			return 0, fmt.Errorf("invalid output parameter type %v", val.Type())
		}

		ctr := 0
		for item := range s.m {
			val.Index(ctr).Set(reflect.ValueOf(item))
			ctr++
		}
		return ctr, nil
	} else {
		return 0, fmt.Errorf("invalid type %v", val.Type())
	}
}
