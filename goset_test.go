package goset

import "testing"

func TestInitial(t *testing.T) {
	NewSet("hello", "world", 1, 1.5, "world")
}

func TestDua(t *testing.T) {
	s := NewStrictSet("hello", "world", "go", "further")
	s.Remove("hello")
}

func TestDistinct(t *testing.T) {
	var (
		out []string
		err error
		len int
	)
	s := NewStrictSet("hello", "world", "go", "further")

	len, err = s.Distinct(&out)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if len != 4 {
		t.Fatalf("unexpected length %v", len)
	}
}
