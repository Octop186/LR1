package main

import (
	"fmt"
	"io"
)

type LNode struct {
	value string
	next  *LNode
	prev  *LNode
}

type List struct {
	head *LNode
	tail *LNode
}

func (l *List) LPushFront(v string) {
	n := &LNode{value: v, next: l.head}
	if l.head == nil {
		l.tail = n
	} else {
		l.head.prev = n
	}
	l.head = n
}

func (l *List) LPushBack(v string) {
	n := &LNode{value: v, prev: l.tail}
	if l.tail == nil {
		l.head = n
	} else {
		l.tail.next = n
	}
	l.tail = n
}

func (l *List) findNode(v string) *LNode {
	for cur := l.head; cur != nil; cur = cur.next {
		if cur.value == v {
			return cur
		}
	}
	return nil
}

func (l *List) LPushAfter(ref, v string) error {
	cur := l.findNode(ref)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	n := &LNode{value: v, prev: cur, next: cur.next}
	if cur.next == nil {
		l.tail = n
	} else {
		cur.next.prev = n
	}
	cur.next = n
	return nil
}

func (l *List) LPushBefore(ref, v string) error {
	cur := l.findNode(ref)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	n := &LNode{value: v, next: cur, prev: cur.prev}
	if cur.prev == nil {
		l.head = n
	} else {
		cur.prev.next = n
	}
	cur.prev = n
	return nil
}

func (l *List) deleteNode(cur *LNode) {
	if cur.prev == nil {
		l.head = cur.next
	} else {
		cur.prev.next = cur.next
	}
	if cur.next == nil {
		l.tail = cur.prev
	} else {
		cur.next.prev = cur.prev
	}
}

func (l *List) LDelFront() error {
	if l.head == nil {
		return fmt.Errorf("ОШИБКА: список пуст")
	}
	l.deleteNode(l.head)
	return nil
}

func (l *List) LDelBack() error {
	if l.tail == nil {
		return fmt.Errorf("ОШИБКА: список пуст")
	}
	l.deleteNode(l.tail)
	return nil
}

func (l *List) LDelAfter(ref string) error {
	cur := l.findNode(ref)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	if cur.next == nil {
		return fmt.Errorf("ОШИБКА: после узла '%s' нет элемента", ref)
	}
	l.deleteNode(cur.next)
	return nil
}

func (l *List) LDelBefore(ref string) error {
	cur := l.findNode(ref)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: узла '%s' нет в списке", ref)
	}
	if cur.prev == nil {
		return fmt.Errorf("ОШИБКА: перед узлом '%s' нет элемента", ref)
	}
	l.deleteNode(cur.prev)
	return nil
}

func (l *List) LDelValue(v string) error {
	cur := l.findNode(v)
	if cur == nil {
		return fmt.Errorf("ОШИБКА: значения '%s' нет в списке", v)
	}
	l.deleteNode(cur)
	return nil
}

func (l *List) LFind(v string) int {
	idx := 0
	for cur := l.head; cur != nil; cur, idx = cur.next, idx+1 {
		if cur.value == v {
			return idx
		}
	}
	return -1
}

func (l *List) LPrint(w io.Writer) {
	fmt.Fprint(w, "L[")
	for cur := l.head; cur != nil; cur = cur.next {
		fmt.Fprint(w, cur.value, " ")
	}
	fmt.Fprintln(w, "]")
}

func (l *List) LPrintReverse(w io.Writer) {
	fmt.Fprint(w, "Lr[")
	for cur := l.tail; cur != nil; cur = cur.prev {
		fmt.Fprint(w, cur.value, " ")
	}
	fmt.Fprintln(w, "]")
}

func (l *List) LPrintRec(w io.Writer) {
	fmt.Fprint(w, "L[")
	var rec func(*LNode)
	rec = func(n *LNode) {
		if n == nil {
			return
		}
		fmt.Fprint(w, n.value, " ")
		rec(n.next)
	}
	rec(l.head)
	fmt.Fprintln(w, "]")
}
