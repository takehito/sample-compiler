package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

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
	reader := []rune(flag.Arg(0))
	numString, err := strtol(&reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Print("	.global main\n")
	fmt.Print("main:\n")
	fmt.Printf("	mov $%d, %%rax\n", numString)
	for len(reader) > 0 {
		char := reader[0]
		reader = reader[1:]
		if char == '+' {
			num, err := strtol(&reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			fmt.Printf("	add $%d, %%rax\n", num)
			continue
		} else if char == '-' {
			i, err := strtol(&reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			fmt.Printf("	sub $%d, %%rax\n", i)
			continue
		} else {
			fmt.Fprintf(os.Stderr, "unexpected character: '%c'\n", char)
			os.Exit(1)
		}
	}
	fmt.Print("	ret\n")
}
