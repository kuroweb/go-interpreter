package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kuromitsu0104/go-interpreter/monkey/02/lexer"
	"github.com/kuromitsu0104/go-interpreter/monkey/02/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
