# Лабораторная работа № 1 — абстрактные структуры данных (Go)

Это многофайловая версия DBMS на Go с тем же интерфейсом, что и предоставленный C++-проект.
Реализованы массив, односвязный и двусвязный списки, стек, очередь, Complete Binary Tree.

## Сборка на Linux

Нужен Go 1.20 или новее.

```bash
cd lab1_go
go build -o dbms .
```

Проверка:

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Примеры

```bash
rm -f file.data
./dbms --file file.data --query 'MPUSH myarray apple'
./dbms --file file.data --query 'MPUSH myarray banana'
./dbms --file file.data --query 'MPUSHIDX myarray 1 cherry'
./dbms --file file.data --query 'PRINT myarray'
# M[apple cherry banana ]
# OK

./dbms --file file.data --query 'QPUSH myqueue first'
./dbms --file file.data --query 'QPUSH myqueue second'
./dbms --file file.data --query 'QGET myqueue'
# first
./dbms --file file.data --query 'QPOP myqueue'
# first
./dbms --file file.data --query 'PRINT myqueue'
# Q[second ]
# OK

./dbms --file file.data --query 'TINSERT mytree 10'
./dbms --file file.data --query 'TINSERT mytree 20'
./dbms --file file.data --query 'TINSERT mytree 30'
./dbms --file file.data --query 'TISCOMPLETE mytree'
# TRUE
./dbms --file file.data --query 'TPRINTTREE mytree'
```

## Команды

| Структура | Поддерживаемые команды |
|---|---|
| Массив M | MPUSH, MPUSHIDX, MGET, MSET, MDEL, MLEN |
| Односвязный список F | FPUSH, FPUSHB, FPUSHAFTER, FPUSHBEFORE, FDEL, FDELBACK, FDELAFTER, FDELBEFORE, FDELVAL, FFIND, FPRINT, FPRINTREC |
| Двусвязный список L | LPUSH, LPUSHB, LPUSHAFTER, LPUSHBEFORE, LDEL, LDELBACK, LDELAFTER, LDELBEFORE, LDELVAL, LFIND, LPRINT, LPRINTR, LPRINTREC |
| Стек S | SPUSH, SPOP, STOP, SPRINT |
| Очередь Q | QPUSH, QPOP, QGET, QPRINT |
| Дерево T | TINSERT, TFIND, TPRINT, TPRINTTREE, TISCOMPLETE, TISFULL |
| Любая структура | PRINT <name> |

Индексация массива начинается с 0. Значения-строки передаются без пробелов.
Для `*AFTER` и `*BEFORE` параметр `ref` — значение узла; используется первое совпадение.
Имена контейнеров могут повторяться у разных типов, но при `PRINT` порядок поиска такой:
массив → односвязный список → двусвязный список → стек → очередь → дерево.

**Файл истории**: сохраняются только успешные команды, меняющие состояние (включая `SPOP` и `QPOP`).
Считывание, вывод, поиск и неудавшиеся изменения в него не записываются. Восстановление
состояния происходит перед каждой командой; старые результаты `PRINT` не дублируются.
Можно начать с чистой истории командой `rm -f file.data`.

Реализация дерева на Go добавлена для полноты; в ТЗ деревья предусмотрены в C++.
Вставка дерева использует ту же схему обхода по индексам, что и исходник C++,
поэтому в худшем случае занимает `O(n)`, а не `O(log n)`; поиск и обход тоже `O(n)`.

## Организация файлов

- `main.go` — аргументы `--file` и `--query`;
- `dbms.go` — хранение структур, обработчик команд, чтение/запись истории;
- `array.go`, `forward_list.go`, `list.go`, `stack.go`, `queue.go`, `tree.go` — реализации АТД;
- `dbms_test.go` — автоматические тесты операций и сохранения данных.
