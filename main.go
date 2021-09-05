package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// 戻り値の第二引数は数値として読み込んだruneの数
func strtol(p io.RuneReader) (int, error) {
	var numString strings.Builder
	for {
		ch, _, err := p.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		if unicode.IsDigit(ch) {
			_, err := numString.WriteRune(ch)
			if err != nil {
				return 0, err
			}
			continue
		}
		break
	}
	num, err := strconv.Atoi(numString.String())
	if err != nil {
		return num, nil
	}
	return num, nil
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid number of arguments\n", os.Args[0])
		os.Exit(1)
	}
	p := flag.Arg(0)
	reader := strings.NewReader(p)
	numString, err := strtol(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Print("	.global main\n")
	fmt.Print("main:\n")
	fmt.Printf("	mov $%d, %%rax\n", numString)
	for {
		ch, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if ch == '+' {
			i, err := strtol(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			fmt.Printf("	add $%d, %%rax\n", i)
			continue
		}
		if ch == '-' {
			i, err := strtol(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			fmt.Printf("	sub $%d, %%rax\n", i)
			continue
		}

		fmt.Fprintf(os.Stderr, "unexpected character: '%c'\n", ch)
		os.Exit(1)
	}
	fmt.Print("	ret\n")
}
