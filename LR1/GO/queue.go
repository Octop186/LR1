package main

import (
	"fmt"
	"io"
)

type Queue struct {
	head *FNode
	tail *FNode
}

func (q *Queue) QEmpty() bool {
	return q.head == nil
}

func (q *Queue) QPush(v string) {
	n := &FNode{value: v}
	if q.tail == nil {
		q.head = n
	} else {
		q.tail.next = n
	}
	q.tail = n
}

func (q *Queue) QPop() (string, error) {
	if q.QEmpty() {
		return "", fmt.Errorf("ОШИБКА: очередь пуста")
	}
	v := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return v, nil
}

func (q *Queue) QFront() (string, error) {
	if q.QEmpty() {
		return "", fmt.Errorf("ОШИБКА: очередь пуста")
	}
	return q.head.value, nil
}

func (q *Queue) QPrint(w io.Writer) {
	fmt.Fprint(w, "Q[")
	for cur := q.head; cur != nil; cur = cur.next {
		fmt.Fprint(w, cur.value, " ")
	}
	fmt.Fprintln(w, "]")
}
