package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	TOKEN_KIND_RESERVED = iota // 記号
	TOKEN_KIND_NUM             // 整数トークン
	TOKEN_KIND_EOF             // 入力の終了を示すトークン
)

type token struct {
	kind int
	next *token
	val  int
	str  []rune
}

func newToken(kind int, cur *token, str []rune) *token {
	t := &token{
		kind: kind,
		str:  str,
	}
	cur.next = t
	return t
}

func tokenize(str []rune) *token {
	var head token
	cur := &head

	for len(str) > 0 {
		c := str[0]
		if unicode.IsSpace(c) {
			str = str[1:]
			continue
		}
		if c == '-' || c == '+' {
			str = str[1:]
			cur = newToken(TOKEN_KIND_RESERVED, cur, []rune{c})
			continue
		}
		if unicode.IsDigit(c) {
			cur = newToken(TOKEN_KIND_NUM, cur, str)
			cur.val, _ = strtol(&str)
			continue
		}

		log.Fatalln("トークナイズできません")
	}

	newToken(TOKEN_KIND_EOF, cur, str)
	return head.next
}

func (t *token) expectNumber() int {
	if t.kind != TOKEN_KIND_NUM {
		log.Fatalln("数値ではありません")
	}
	val := t.val
	*t = *t.next
	return val
}

// 次のトークンが期待している記号の時はトークンを一つ進めて
// 真を返す。それ以外の場合には偽を返す
func (t *token) consume(r rune) bool {
	if t.kind != TOKEN_KIND_RESERVED || t.str[0] != r {
		return false
	}
	*t = *t.next
	return true
}

func (t *token) expect(r rune) {
	if t.kind != TOKEN_KIND_RESERVED || t.str[0] != r {
		log.Fatalf("'%#U'ではありません, '%#U'\n", r, t.str[0])
	}
	*t = *t.next
}

func (t token) atEOF() bool {
	return t.kind == TOKEN_KIND_EOF
}

// 戻り値の第二引数は数値として読み込んだruneの数
func strtol(r *[]rune) (int, error) {
	var numString strings.Builder
	c := *r
	for len(c) > 0 {
		if unicode.IsDigit(c[0]) {
			if _, err := numString.WriteRune(c[0]); err != nil {
				return 0, err
			}
			c = c[1:]
			continue
		}

		break
	}
	*r = c
	return strconv.Atoi(numString.String())
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid number of arguments\n", os.Args[0])
		os.Exit(1)
	}

	t := tokenize([]rune(flag.Arg(0)))

	//fmt.Print(".intel_syntax noprefix\n")
	fmt.Print("	.global main\n")
	fmt.Print("main:\n")

	fmt.Printf("	mov $%d, %%rax\n", t.expectNumber())

	for !t.atEOF() {
		if t.consume('+') {
			fmt.Printf("	add $%d, %%rax\n", t.expectNumber())
			continue
		}
		t.expect('-')
		fmt.Printf("	sub $%d, %%rax\n", t.expectNumber())
	}
	fmt.Print("	ret\n")
}
