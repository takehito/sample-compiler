package main

import (
	"flag"
	"fmt"
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

const (
	NODE_KIND_ADD = iota
	NODE_KIND_SUB
	NODE_KIND_MUL
	NODE_KIND_DIV
	NODE_KIND_NUM
)

type node struct {
	nodeKind int
	lhs      *node
	rhs      *node
	val      int
}

func newNode(nodeKind int, lhs *node, rhs *node) *node {
	return &node{
		nodeKind: nodeKind,
		lhs:      lhs,
		rhs:      rhs,
	}
}

func newNodeNum(val int) *node {
	return &node{
		nodeKind: NODE_KIND_NUM,
		val:      val,
	}
}

func (n *node) gen() {
	if n.nodeKind == NODE_KIND_NUM {
		fmt.Printf("	push %d\n", n.val)
		return
	}

	n.lhs.gen()
	n.rhs.gen()

	fmt.Printf("	pop rdi\n")
	fmt.Printf("	pop rax\n")

	switch n.nodeKind {
	case NODE_KIND_ADD:
		fmt.Printf("	add rax, rdi\n")
	case NODE_KIND_SUB:
		fmt.Printf("	sub rax, rdi\n")
	case NODE_KIND_MUL:
		fmt.Printf("	imul rax, rdi\n")
	case NODE_KIND_DIV:
		fmt.Printf("	cqo\n")
		fmt.Printf("	idiv rdi\n")
	}

	fmt.Printf("	push rax\n")
}

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
		if c == '-' || c == '+' || c == '*' || c == '/' || c == '(' || c == ')' {
			str = str[1:]
			cur = newToken(TOKEN_KIND_RESERVED, cur, []rune{c})
			continue
		}
		if unicode.IsDigit(c) {
			cur = newToken(TOKEN_KIND_NUM, cur, str)
			cur.val, _ = strtol(&str)
			continue
		}

		errorAt(str, "トークナイズできません")
	}

	newToken(TOKEN_KIND_EOF, cur, str)
	return head.next
}

func (t *token) expr() *node {
	n := t.mul()

	for {
		if t.consume('+') {
			n = newNode(NODE_KIND_ADD, n, t.mul())
		} else if t.consume('-') {
			n = newNode(NODE_KIND_SUB, n, t.mul())
		} else {
			return n
		}
	}
}

func (t *token) mul() *node {
	n := t.unary()

	for {
		if t.consume('*') {
			n = newNode(NODE_KIND_MUL, n, t.unary())
		} else if t.consume('/') {
			n = newNode(NODE_KIND_DIV, n, t.unary())
		} else {
			return n
		}
	}
}

func (t *token) unary() *node {
	if t.consume('+') {
		return t.expr()
	}
	if t.consume('-') {
		return newNode(NODE_KIND_SUB, newNodeNum(0), t.primary())
	}
	return t.primary()
}

func (t *token) primary() *node {
	// 次のトークンが"("なら"("+expr+")"のはず
	if t.consume('(') {
		n := t.expr()
		t.expect(')')
		return n
	}

	// そうでなければ通知のはず
	return newNodeNum(t.expectNumber())
}

func (t *token) expectNumber() int {
	if t.kind != TOKEN_KIND_NUM {
		errorAt(t.str, "数値ではありません")
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
		errorAt(t.str, "'%#U'ではありません, '%#U'\n", r, t.str[0])
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

var userInput string

func errorAt(loc []rune, str string, form ...interface{}) {
	pos := len(userInput) - len(loc)
	fmt.Fprintf(os.Stderr, "%s\n", userInput)
	fmt.Fprintf(os.Stderr, "%*s", pos, " ")
	fmt.Fprint(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, str, form...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid number of arguments\n", os.Args[0])
		os.Exit(1)
	}

	userInput = flag.Arg(0)
	t := tokenize([]rune(userInput))
	n := t.expr()

	fmt.Print(".intel_syntax noprefix\n")
	fmt.Print("	.global main\n")
	fmt.Print("main:\n")

	n.gen()

	fmt.Print("	pop rax\n")
	fmt.Print("	ret\n")
}
