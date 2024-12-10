package main

import (
	"fmt"
	"go_interpreter/evaluator"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"go_interpreter/repl"
	"io"
	"os"
	"os/user"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) == 2 {

		// f, err := os.Open(os.Args[1])
		// check(err)

		text, err := os.ReadFile(os.Args[1])
		check(err)
		env := object.NewEnvironment()
		line := string(text)
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
		}
		// io.WriteString(os.Stdout, program.String())
		// io.WriteString(os.Stdout, "\n")

		evaluated := evaluator.Eval(program, env)
		if evaluated.Type() == object.ERROR_OBJ {
			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	} else {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Println(` _._     _,-'""` + "`" + `-._` + "\n" +
			`(,-.` + "`" + `._,'(       |\` + "`" + `-/|` + "\n" +
			`    ` + "`" + `-.-' \ )-` + "`" + `( , o o)` + "\n" +
			`          ` + "`" + `-    \` + "`" + `_` + "`" + ` '-`)

		fmt.Printf("\nHello %s! This is the meatball programming language\n\n",
			user.Username)
		repl.Start(os.Stdin, os.Stdout)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
