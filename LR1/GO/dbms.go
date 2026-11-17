package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DB struct {
	arrays map[string]*Array
	flists map[string]*ForwardList
	lists  map[string]*List
	stacks map[string]*Stack
	queues map[string]*Queue
	trees  map[string]*CompleteBinaryTree
}

func NewDB() *DB {
	return &DB{
		arrays: make(map[string]*Array),
		flists: make(map[string]*ForwardList),
		lists:  make(map[string]*List),
		stacks: make(map[string]*Stack),
		queues: make(map[string]*Queue),
		trees:  make(map[string]*CompleteBinaryTree),
	}
}

func (db *DB) array(name string) *Array {
	if db.arrays[name] == nil {
		db.arrays[name] = &Array{}
	}
	return db.arrays[name]
}

func (db *DB) flist(name string) *ForwardList {
	if db.flists[name] == nil {
		db.flists[name] = &ForwardList{}
	}
	return db.flists[name]
}

func (db *DB) list(name string) *List {
	if db.lists[name] == nil {
		db.lists[name] = &List{}
	}
	return db.lists[name]
}

func (db *DB) stack(name string) *Stack {
	if db.stacks[name] == nil {
		db.stacks[name] = &Stack{}
	}
	return db.stacks[name]
}

func (db *DB) queue(name string) *Queue {
	if db.queues[name] == nil {
		db.queues[name] = &Queue{}
	}
	return db.queues[name]
}

func (db *DB) tree(name string) *CompleteBinaryTree {
	if db.trees[name] == nil {
		db.trees[name] = &CompleteBinaryTree{}
	}
	return db.trees[name]
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("ОШИБКА: значение — не число")
	}
	start := 0
	if s[0] == '-' {
		start = 1
	}
	if start == len(s) {
		return 0, fmt.Errorf("ОШИБКА: значение — не число")
	}
	for _, c := range s[start:] {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("ОШИБКА: значение — не число")
		}
	}
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("ОШИБКА: число выходит за пределы int")
	}
	return int(n), nil
}

func isError(res string) bool {
	return strings.HasPrefix(res, "ОШИБКА")
}

func commandError(op, usage string) (string, bool) {
	return "ОШИБКА: " + op + " " + usage, false
}

func result(err error) (string, bool) {
	if err != nil {
		return err.Error(), false
	}
	return "OK", true
}

// Exec выполняет одну команду и возвращает (ответ, изменены ли данные).
// Побочный вывод (например, PRINT) передаётся в out.
func Exec(db *DB, line string, out io.Writer) (string, bool) {
	t := strings.Fields(line)
	if len(t) == 0 {
		return "ОШИБКА: пустая команда", false
	}
	op := t[0]
	if op == "PRINT" {
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		name := t[1]
		if v := db.arrays[name]; v != nil {
			v.MPrint(out)
			return "OK", false
		}
		if v := db.flists[name]; v != nil {
			v.FPrint(out)
			return "OK", false
		}
		if v := db.lists[name]; v != nil {
			v.LPrint(out)
			return "OK", false
		}
		if v := db.stacks[name]; v != nil {
			v.SPrint(out)
			return "OK", false
		}
		if v := db.queues[name]; v != nil {
			v.QPrint(out)
			return "OK", false
		}
		if v := db.trees[name]; v != nil {
			v.TPrint(out)
			return "OK", false
		}
		return fmt.Sprintf("ОШИБКА: контейнер '%s' не найден", name), false
	}

	switch op {
	case "MPUSH":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.array(t[1]).MPush(t[2])
		return t[2], true
	case "MPUSHIDX":
		if len(t) != 4 {
			return commandError(op, "<name> <idx> <value>")
		}
		idx, err := parseInt(t[2])
		if err != nil {
			return "ОШИБКА: idx — не число", false
		}
		if err := db.array(t[1]).MPushIdx(idx, t[3]); err != nil {
			return err.Error(), false
		}
		return t[3], true
	case "MGET":
		if len(t) != 3 {
			return commandError(op, "<name> <idx>")
		}
		idx, err := parseInt(t[2])
		if err != nil {
			return "ОШИБКА: idx — не число", false
		}
		v, err := db.array(t[1]).MGet(idx)
		if err != nil {
			return err.Error(), false
		}
		return v, false
	case "MDEL":
		if len(t) != 3 {
			return commandError(op, "<name> <idx>")
		}
		idx, err := parseInt(t[2])
		if err != nil {
			return "ОШИБКА: idx — не число", false
		}
		return result(db.array(t[1]).MDelInd(idx))
	case "MSET":
		if len(t) != 4 {
			return commandError(op, "<name> <idx> <value>")
		}
		idx, err := parseInt(t[2])
		if err != nil {
			return "ОШИБКА: idx — не число", false
		}
		if err := db.array(t[1]).MSetInd(idx, t[3]); err != nil {
			return err.Error(), false
		}
		return t[3], true
	case "MLEN":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		return strconv.Itoa(db.array(t[1]).MLen()), false

	case "FPUSH":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.flist(t[1]).FPushFront(t[2])
		return t[2], true
	case "FPUSHB":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.flist(t[1]).FPushBack(t[2])
		return t[2], true
	case "FPUSHAFTER":
		if len(t) != 4 {
			return commandError(op, "<name> <ref> <value>")
		}
		if err := db.flist(t[1]).FPushAfter(t[2], t[3]); err != nil {
			return err.Error(), false
		}
		return t[3], true
	case "FPUSHBEFORE":
		if len(t) != 4 {
			return commandError(op, "<name> <ref> <value>")
		}
		if err := db.flist(t[1]).FPushBefore(t[2], t[3]); err != nil {
			return err.Error(), false
		}
		return t[3], true
	case "FDEL":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		return result(db.flist(t[1]).FDelFront())
	case "FDELBACK":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		return result(db.flist(t[1]).FDelBack())
	case "FDELAFTER":
		if len(t) != 3 {
			return commandError(op, "<name> <ref>")
		}
		return result(db.flist(t[1]).FDelAfter(t[2]))
	case "FDELBEFORE":
		if len(t) != 3 {
			return commandError(op, "<name> <ref>")
		}
		return result(db.flist(t[1]).FDelBefore(t[2]))
	case "FDELVAL":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		return result(db.flist(t[1]).FDelValue(t[2]))
	case "FFIND":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		idx := db.flist(t[1]).FFind(t[2])
		if idx >= 0 {
			fmt.Fprintf(out, "Элемент был найден на позиции %d\n", idx)
		} else {
			fmt.Fprintln(out, "Элемент не был найден в списке")
		}
		return "OK", false
	case "FPRINT":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.flist(t[1]).FPrint(out)
		return "OK", false
	case "FPRINTREC":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.flist(t[1]).FPrintRec(out)
		return "OK", false

	case "LPUSH":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.list(t[1]).LPushFront(t[2])
		return t[2], true
	case "LPUSHB":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.list(t[1]).LPushBack(t[2])
		return t[2], true
	case "LPUSHAFTER":
		if len(t) != 4 {
			return commandError(op, "<name> <ref> <value>")
		}
		if err := db.list(t[1]).LPushAfter(t[2], t[3]); err != nil {
			return err.Error(), false
		}
		return t[3], true
	case "LPUSHBEFORE":
		if len(t) != 4 {
			return commandError(op, "<name> <ref> <value>")
		}
		if err := db.list(t[1]).LPushBefore(t[2], t[3]); err != nil {
			return err.Error(), false
		}
		return t[3], true
	case "LDEL":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		return result(db.list(t[1]).LDelFront())
	case "LDELBACK":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		return result(db.list(t[1]).LDelBack())
	case "LDELAFTER":
		if len(t) != 3 {
			return commandError(op, "<name> <ref>")
		}
		return result(db.list(t[1]).LDelAfter(t[2]))
	case "LDELBEFORE":
		if len(t) != 3 {
			return commandError(op, "<name> <ref>")
		}
		return result(db.list(t[1]).LDelBefore(t[2]))
	case "LDELVAL":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		return result(db.list(t[1]).LDelValue(t[2]))
	case "LFIND":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		idx := db.list(t[1]).LFind(t[2])
		if idx >= 0 {
			fmt.Fprintf(out, "Элемент был найден на позиции %d\n", idx)
		} else {
			fmt.Fprintln(out, "Элемент не был найден в списке")
		}
		return "OK", false
	case "LPRINT":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.list(t[1]).LPrint(out)
		return "OK", false
	case "LPRINTR":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.list(t[1]).LPrintReverse(out)
		return "OK", false
	case "LPRINTREC":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.list(t[1]).LPrintRec(out)
		return "OK", false

	case "SPUSH":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.stack(t[1]).SPush(t[2])
		return t[2], true
	case "SPOP":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		v, err := db.stack(t[1]).SPop()
		if err != nil {
			return err.Error(), false
		}
		return v, true
	case "STOP":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		v, err := db.stack(t[1]).STop()
		if err != nil {
			return err.Error(), false
		}
		return v, false
	case "SPRINT":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.stack(t[1]).SPrint(out)
		return "OK", false

	case "QPUSH":
		if len(t) != 3 {
			return commandError(op, "<name> <value>")
		}
		db.queue(t[1]).QPush(t[2])
		return t[2], true
	case "QPOP":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		v, err := db.queue(t[1]).QPop()
		if err != nil {
			return err.Error(), false
		}
		return v, true
	case "QGET":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		v, err := db.queue(t[1]).QFront()
		if err != nil {
			return err.Error(), false
		}
		return v, false
	case "QPRINT":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.queue(t[1]).QPrint(out)
		return "OK", false

	case "TINSERT":
		if len(t) != 3 {
			return commandError(op, "<name> <int>")
		}
		v, err := parseInt(t[2])
		if err != nil {
			return err.Error(), false
		}
		db.tree(t[1]).TInsert(v)
		return t[2], true
	case "TFIND":
		if len(t) != 3 {
			return commandError(op, "<name> <int>")
		}
		v, err := parseInt(t[2])
		if err != nil {
			return err.Error(), false
		}
		d, ok := db.tree(t[1]).TFind(v)
		if ok {
			fmt.Fprintf(out, "Элемент был найден на глубине %d\n", d)
		} else {
			fmt.Fprintln(out, "Элемент не был найден в дереве")
		}
		return "OK", false
	case "TPRINT":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.tree(t[1]).TPrint(out)
		return "OK", false
	case "TPRINTTREE":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		db.tree(t[1]).TPrintTree(out)
		return "OK", false
	case "TISCOMPLETE":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		if db.tree(t[1]).TIsComplete() {
			return "TRUE", false
		}
		return "FALSE", false
	case "TISFULL":
		if len(t) != 2 {
			return commandError(op, "<name>")
		}
		if db.tree(t[1]).TIsFull() {
			return "TRUE", false
		}
		return "FALSE", false
	}
	return fmt.Sprintf("ОШИБКА: неизвестная команда '%s'", op), false
}

// Load восстанавливает состояние по успешным изменяющим командам из файла.
func Load(db *DB, path string) error {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		Exec(db, line, io.Discard)
	}
	return scanner.Err()
}

// Save добавляет только команды, изменяющие данные.
func Save(path, command string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, command)
	return err
}
