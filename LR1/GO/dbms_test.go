package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func check(t *testing.T, db *DB, cmd, want string, changed bool) string {
	t.Helper()
	res, mutated := Exec(db, cmd, io.Discard)
	if res != want || mutated != changed {
		t.Fatalf("%q: ответ=%q changed=%v, ожидали %q changed=%v", cmd, res, mutated, want, changed)
	}
	return res
}

func checkPrint(t *testing.T, db *DB, cmd, want string) {
	t.Helper()
	var buf bytes.Buffer
	res, changed := Exec(db, cmd, &buf)
	if res != "OK" || changed || buf.String() != want {
		t.Fatalf("%q: ответ=%q, changed=%v, вывод=%q, ожидали %q", cmd, res, changed, buf.String(), want)
	}
}

func TestArray(t *testing.T) {
	db := NewDB()
	check(t, db, "MPUSH a apple", "apple", true)
	check(t, db, "MPUSH a banana", "banana", true)
	check(t, db, "MPUSHIDX a 1 cherry", "cherry", true)
	checkPrint(t, db, "PRINT a", "M[apple cherry banana ]\n")
	check(t, db, "MGET a 1", "cherry", false)
	check(t, db, "MSET a 0 orange", "orange", true)
	check(t, db, "MDEL a 1", "OK", true)
	check(t, db, "MLEN a", "2", false)
	checkPrint(t, db, "PRINT a", "M[orange banana ]\n")
	check(t, db, "MDEL a 20", "ОШИБКА: индекс 20 не существует", false)
	check(t, db, "MGET a -1", "ОШИБКА: индекс -1 не существует", false)
	check(t, db, "MPUSHIDX a 4 x", "ОШИБКА: индекс 4 не существует", false)
	check(t, db, "MSET a 3 x", "ОШИБКА: индекс 3 не существует", false)
	check(t, db, "MGET a 999999999999999999999", "ОШИБКА: idx — не число", false)
}

func TestForwardList(t *testing.T) {
	db := NewDB()
	check(t, db, "FPUSH f b", "b", true)
	check(t, db, "FPUSH f a", "a", true)
	check(t, db, "FPUSHB f d", "d", true)
	check(t, db, "FPUSHAFTER f b c", "c", true)
	check(t, db, "FPUSHBEFORE f a z", "z", true)
	checkPrint(t, db, "FPRINT f", "F[z a b c d ]\n")
	checkPrint(t, db, "FPRINTREC f", "F[z a b c d ]\n")
	check(t, db, "FDELBEFORE f a", "OK", true)
	check(t, db, "FDELAFTER f b", "OK", true)
	check(t, db, "FDELVAL f b", "OK", true)
	checkPrint(t, db, "PRINT f", "F[a d ]\n")
	check(t, db, "FDELAFTER f d", "ОШИБКА: после узла 'd' нет элемента", false)
	check(t, db, "FDELBEFORE f a", "ОШИБКА: перед головой нет элемента", false)
	check(t, db, "FDELVAL f none", "ОШИБКА: значения 'none' нет в списке", false)
	check(t, db, "FDEL f", "OK", true)
	check(t, db, "FDELBACK f", "OK", true)
	check(t, db, "FDEL f", "ОШИБКА: список пуст", false)
}

func TestDoublyList(t *testing.T) {
	db := NewDB()
	check(t, db, "LPUSHB l b", "b", true)
	check(t, db, "LPUSH l a", "a", true)
	check(t, db, "LPUSHAFTER l b d", "d", true)
	check(t, db, "LPUSHBEFORE l d c", "c", true)
	checkPrint(t, db, "LPRINT l", "L[a b c d ]\n")
	checkPrint(t, db, "LPRINTR l", "Lr[d c b a ]\n")
	check(t, db, "LDELBEFORE l d", "OK", true)
	check(t, db, "LDELAFTER l b", "OK", true)
	checkPrint(t, db, "LPRINTR l", "Lr[b a ]\n")
	check(t, db, "LDELVAL l a", "OK", true)
	check(t, db, "LDEL l", "OK", true)
	check(t, db, "LDELBACK l", "ОШИБКА: список пуст", false)
	check(t, db, "LPUSHB l final", "final", true)
	checkPrint(t, db, "LPRINTR l", "Lr[final ]\n")
}

func TestStackQueue(t *testing.T) {
	db := NewDB()
	check(t, db, "SPUSH s one", "one", true)
	check(t, db, "SPUSH s two", "two", true)
	check(t, db, "STOP s", "two", false)
	check(t, db, "SPOP s", "two", true)
	check(t, db, "SPOP s", "one", true)
	check(t, db, "SPOP s", "ОШИБКА: стек пуст", false)
	check(t, db, "QPUSH q first", "first", true)
	check(t, db, "QPUSH q second", "second", true)
	check(t, db, "QGET q", "first", false)
	check(t, db, "QPOP q", "first", true)
	check(t, db, "QPOP q", "second", true)
	check(t, db, "QPOP q", "ОШИБКА: очередь пуста", false)
	check(t, db, "QPUSH q third", "third", true)
	checkPrint(t, db, "PRINT q", "Q[third ]\n")
}

func TestTree(t *testing.T) {
	db := NewDB()
	for _, cmd := range []string{"TINSERT t 10", "TINSERT t 20", "TINSERT t 30", "TINSERT t 40", "TINSERT t 50"} {
		res, changed := Exec(db, cmd, io.Discard)
		if !changed || strings.HasPrefix(res, "ОШИБКА") {
			t.Fatalf("%s: %s", cmd, res)
		}
	}
	check(t, db, "TISCOMPLETE t", "TRUE", false)
	check(t, db, "TISFULL t", "TRUE", false)
	checkPrint(t, db, "TPRINT t", "T[40 20 50 10 30 ]\n")
	check(t, db, "TINSERT t 60", "60", true)
	check(t, db, "TISFULL t", "FALSE", false)
	var b bytes.Buffer
	res, changed := Exec(db, "TFIND t 50", &b)
	if res != "OK" || changed || b.String() != "Элемент был найден на глубине 2\n" {
		t.Fatalf("TFIND: %q %v %q", res, changed, b.String())
	}
	check(t, db, "TINSERT t 2147483648", "ОШИБКА: число выходит за пределы int", false)
}

func TestReplayOnlyMutations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.data")
	db := NewDB()
	commands := []string{"QPUSH q first", "QPUSH q second", "QGET q", "QPOP q", "QPRINT q", "QPOP q", "QPOP q", "QPUSH q after"}
	saved := 0
	for _, cmd := range commands {
		res, changed := Exec(db, cmd, io.Discard)
		if changed {
			if err := Save(path, cmd); err != nil {
				t.Fatal(err)
			}
			saved++
		}
		if cmd == "QPOP q" && res == "" {
			t.Fatal("пустой ответ")
		}
	}
	if saved != 5 {
		t.Fatalf("сохранено команд: %d", saved)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "QGET") || strings.Contains(string(content), "QPRINT") {
		t.Fatalf("команды чтения сохранены: %q", content)
	}
	restored := NewDB()
	if err := Load(restored, path); err != nil {
		t.Fatal(err)
	}
	checkPrint(t, restored, "PRINT q", "Q[after ]\n")
}
