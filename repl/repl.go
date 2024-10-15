package repl

import (
	"bufio"
	"fmt"
	"go_interpreter/evaluator"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"io"
)

const PROMPT = "｡＾･ｪ･＾｡ >> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

// func Start(in io.Reader, out io.Writer) {
// 	scanner := bufio.NewScanner(in)
// 	for {
// 		fmt.Printf(PROMPT)
// 		scanned := scanner.Scan()
// 		if !scanned {
// 			return
// 		}
// 		line := scanner.Text()
// 		l := lexer.New(line)
// 		p := parser.New(l)
// 		program := p.ParseProgram()
// 		if len(p.Errors()) != 0 {
// 			printParserErrors(out, p.Errors())
// 			continue
// 		}
// 		io.WriteString(out, program.String())
// 		io.WriteString(out, "\n")
// 	}
// }
//
// func printParserErrors(out io.Writer, errors []string) {
// 	for _, msg := range errors {
// 		io.WriteString(out, "\t"+msg+"\n")
// 	}
// }
