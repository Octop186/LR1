package main

import (
	"fmt"
	"io"
)

type Stack struct {
	head *FNode
}

func (s *Stack) SEmpty() bool {
	return s.head == nil
}

func (s *Stack) SPush(v string) {
	s.head = &FNode{value: v, next: s.head}
}

func (s *Stack) SPop() (string, error) {
	if s.SEmpty() {
		return "", fmt.Errorf("ОШИБКА: стек пуст")
	}
	v := s.head.value
	s.head = s.head.next
	return v, nil
}

func (s *Stack) STop() (string, error) {
	if s.SEmpty() {
		return "", fmt.Errorf("ОШИБКА: стек пуст")
	}
	return s.head.value, nil
}

func (s *Stack) SPrint(w io.Writer) {
	fmt.Fprint(w, "S[")
	for cur := s.head; cur != nil; cur = cur.next {
		fmt.Fprint(w, cur.value, " ")
	}
	fmt.Fprintln(w, "]")
}
