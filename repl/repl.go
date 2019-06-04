package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/user/courseWork/ast"
	"github.com/user/courseWork/evaluator"
	"github.com/user/courseWork/lexer"
	"github.com/user/courseWork/object"
	"github.com/user/courseWork/parser"
	"io"
	"sync"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {

	wg := new(sync.WaitGroup)
	wg.Add(3)

	analysis := make(chan string)
	evaluation := make(chan *ast.Program)
	printing := make(chan string)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(in)
		for {
			fmt.Print(PROMT)
			scanned := scanner.Scan()
			if !scanned {
				close(analysis)
				return
			}
			line := scanner.Text()
			analysis <- line
			ostr := <-printing
			if ostr != "" {
				io.WriteString(out, ostr)
				io.WriteString(out, "\n")
			}
		}
	}()

	go func() {
		defer wg.Done()
		for line := range analysis {
			l := lexer.New(line)
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				var err bytes.Buffer
				for _, msg := range p.Errors() {
					err.WriteString("\t" + msg + "\n")
				}
				printing <- err.String()
				continue
			}
			evaluation <- program
		}
		close(evaluation)
	}()

	go func() {
		defer wg.Done()
		env := object.NewEnvironment()
		for program := range evaluation {
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				printing <- evaluated.Inspect()
			} else {
				printing <- ""
			}
		}
		close(printing)
	}()

	wg.Wait()
}
