package main

import (
	"fmt"
	"io"
	"strings"
)

type Array struct {
	data []string
}

func (a *Array) MPush(v string) {
	a.data = append(a.data, v)
}

func (a *Array) MPushIdx(idx int, v string) error {
	if idx < 0 || idx > len(a.data) {
		return fmt.Errorf("ОШИБКА: индекс %d не существует", idx)
	}
	a.data = append(a.data, "")
	copy(a.data[idx+1:], a.data[idx:len(a.data)-1])
	a.data[idx] = v
	return nil
}

func (a *Array) MGet(idx int) (string, error) {
	if idx < 0 || idx >= len(a.data) {
		return "", fmt.Errorf("ОШИБКА: индекс %d не существует", idx)
	}
	return a.data[idx], nil
}

func (a *Array) MDelInd(idx int) error {
	if idx < 0 || idx >= len(a.data) {
		return fmt.Errorf("ОШИБКА: индекс %d не существует", idx)
	}
	copy(a.data[idx:], a.data[idx+1:])
	a.data = a.data[:len(a.data)-1]
	return nil
}

func (a *Array) MSetInd(idx int, v string) error {
	if idx < 0 || idx >= len(a.data) {
		return fmt.Errorf("ОШИБКА: индекс %d не существует", idx)
	}
	a.data[idx] = v
	return nil
}

func (a *Array) MLen() int {
	return len(a.data)
}

func (a *Array) MPrint(w io.Writer) {
	fmt.Fprintf(w, "M[%s]\n", withTrailingSpace(a.data))
}

func withTrailingSpace(v []string) string {
	if len(v) == 0 {
		return ""
	}
	return strings.Join(v, " ") + " "
}
