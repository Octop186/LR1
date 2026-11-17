package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	file := flag.String("file", "", "файл данных")
	query := flag.String("query", "", "команда DBMS")
	flag.Parse()

	db := NewDB()
	if *file != "" {
		if err := Load(db, *file); err != nil {
			fmt.Fprintln(os.Stderr, "ОШИБКА: чтение файла:", err)
			os.Exit(1)
		}
	}
	if *query == "" {
		return
	}

	res, changed := Exec(db, *query, os.Stdout)
	fmt.Println(res)
	if changed && *file != "" {
		if err := Save(*file, *query); err != nil {
			fmt.Fprintln(os.Stderr, "ОШИБКА: запись файла:", err)
			os.Exit(1)
		}
	}
}
