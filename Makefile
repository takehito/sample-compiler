chibicc: *.go
	go build -o chibicc

test: chibicc
	./test.sh

clean:
	rm -f chibicc tmp*

.PHONY: test clean 
