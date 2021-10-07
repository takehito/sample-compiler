chibicc: *.go
	go build -o chibicc

test: chibicc
	go test
	./test.sh

clean:
	rm -f chibicc tmp*

.PHONY: test clean 
