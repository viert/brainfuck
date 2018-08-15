package brainfuck

import (
	"bufio"
	"fmt"
	"os"
)

type stack []int

func (s *stack) Push(el int) {
	*s = append(*s, el)
}

func (s *stack) Pop() int {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

type Interpreter struct {
	code   []byte
	memory []byte
	stack
}

func New(code []byte) *Interpreter {
	return &Interpreter{
		code,
		make([]byte, 1, 30000),
		make(stack, 0, 1000),
	}
}

func (i *Interpreter) Run() {
	pc := 0
	dp := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		switch i.code[pc] {
		case '>':
			dp++
			if dp > len(i.memory)-1 {
				i.memory = append(i.memory, 0)
			}
		case '<':
			dp--
			if dp < 0 {
				panic("Data pointer < 0")
			}
		case '+':
			i.memory[dp]++
		case '-':
			i.memory[dp]--
		case '[':
			if i.memory[dp] == 0 {
				oldpc := pc
				braces := 1
				for {
					pc++
					if pc > len(i.code)-1 {
						panic(fmt.Sprintf("No closing brace for [ at %d", oldpc))
					}
					if i.code[pc] == '[' {
						braces++
					}
					if i.code[pc] == ']' {
						braces--
						if braces == 0 {
							break
						}
					}
				}
			} else {
				i.stack.Push(pc - 1)
			}
		case ']':
			pc = i.stack.Pop()
		case '.':
			fmt.Print(string(i.memory[dp]))
		case ',':
			input, _ := reader.ReadString('\n')
			i.memory[dp] = input[0]
		case '\n':
		default:
			panic(fmt.Sprintf("Incorect symbol '%s' at %d", i.code[pc], pc))
		}

		pc++
		if pc == len(i.code) {
			break
		}
	}
}
