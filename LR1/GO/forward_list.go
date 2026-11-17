package main

import (
	"fmt"
	"io"
)

type FNode struct {
	value string
	next  *FNode
}

type ForwardList struct {
	head *FNode
}

func (f *ForwardList) FPushFront(v string) {
	f.head = &FNode{value: v, next: f.head}
}

func (f *ForwardList) FPushBack(v string) {
	n := &FNode{value: v}
	if f.head == nil {
		f.head = n
		return
	}
	cur := f.head
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = n
}

func (f *ForwardList) findNode(v string) *FNode {
	for cur := f.head; cur != nil; cur = cur.next {
		if cur.value == v {
			return cur
		}
	}
	return nil
}

func (f *ForwardList) FPushAfter(ref, v string) error {
	cur := f.findNode(ref)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	cur.next = &FNode{value: v, next: cur.next}
	return nil
}

func (f *ForwardList) FPushBefore(ref, v string) error {
	if f.head == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	if f.head.value == ref {
		f.FPushFront(v)
		return nil
	}
	for cur := f.head; cur.next != nil; cur = cur.next {
		if cur.next.value == ref {
			cur.next = &FNode{value: v, next: cur.next}
			return nil
		}
	}
	return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
}

func (f *ForwardList) FDelFront() error {
	if f.head == nil {
		return fmt.Errorf("ОШИБКА: список пуст")
	}
	f.head = f.head.next
	return nil
}

func (f *ForwardList) FDelBack() error {
	if f.head == nil {
		return fmt.Errorf("ОШИБКА: список пуст")
	}
	if f.head.next == nil {
		f.head = nil
		return nil
	}
	cur := f.head
	for cur.next.next != nil {
		cur = cur.next
	}
	cur.next = nil
	return nil
}

func (f *ForwardList) FDelAfter(ref string) error {
	cur := f.findNode(ref)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	if cur.next == nil {
		return fmt.Errorf("ОШИБКА: после узла '%s' нет элемента", ref)
	}
	cur.next = cur.next.next
	return nil
}

func (f *ForwardList) FDelBefore(ref string) error {
	if f.head == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	if f.head.value == ref {
		return fmt.Errorf("ОШИБКА: перед головой нет элемента")
	}
	if f.head.next != nil && f.head.next.value == ref {
		f.head = f.head.next
		return nil
	}
	for cur := f.head; cur.next != nil && cur.next.next != nil; cur = cur.next {
		if cur.next.next.value == ref {
			cur.next = cur.next.next
			return nil
		}
	}
	return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
}

func (f *ForwardList) FDelValue(v string) error {
	if f.head == nil {
		return fmt.Errorf("ОШИБКА: значения '%s' нет в списке", v)
	}
	if f.head.value == v {
		f.head = f.head.next
		return nil
	}
	for cur := f.head; cur.next != nil; cur = cur.next {
		if cur.next.value == v {
			cur.next = cur.next.next
			return nil
		}
	}
	return fmt.Errorf("ОШИБКА: значения '%s' нет в списке", v)
}

func (f *ForwardList) FFind(v string) int {
	idx := 0
	for cur := f.head; cur != nil; cur, idx = cur.next, idx+1 {
		if cur.value == v {
			return idx
		}
	}
	return -1
}

func (f *ForwardList) FPrint(w io.Writer) {
	fmt.Fprint(w, "F[")
	for cur := f.head; cur != nil; cur = cur.next {
		fmt.Fprint(w, cur.value, " ")
	}
	fmt.Fprintln(w, "]")
}

func (f *ForwardList) FPrintRec(w io.Writer) {
	fmt.Fprint(w, "F[")
	var rec func(*FNode)
	rec = func(n *FNode) {
		if n == nil {
			return
		}
		fmt.Fprint(w, n.value, " ")
		rec(n.next)
	}
	rec(f.head)
	fmt.Fprintln(w, "]")
}
