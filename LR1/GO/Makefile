.PHONY: all test vet clean

all:
	go build -o dbms .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f dbms file.data
