#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
go build -o dbms .
tmpfile="$(mktemp)"
trap 'rm -f "$tmpfile"' EXIT
rm -f "$tmpfile"
run() { ./dbms --file "$tmpfile" --query "$1"; }
[[ "$(run 'MPUSH a apple')" == 'apple' ]]
[[ "$(run 'MPUSH a banana')" == 'banana' ]]
[[ "$(run 'MPUSHIDX a 1 cherry')" == 'cherry' ]]
[[ "$(run 'PRINT a')" == $'M[apple cherry banana ]\nOK' ]]
[[ "$(run 'MDEL a 99')" == 'ОШИБКА: индекс 99 не существует' ]]
[[ "$(run 'QPUSH q first')" == 'first' ]]
[[ "$(run 'QPUSH q second')" == 'second' ]]
[[ "$(run 'QPOP q')" == 'first' ]]
[[ "$(run 'PRINT q')" == $'Q[second ]\nOK' ]]
[[ "$(grep -c '^PRINT\|^MDEL a 99' "$tmpfile" || true)" == '0' ]]
echo 'CLI-проверки Go: OK'
