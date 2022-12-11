package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	answer, err := solveNoSpaceLeftOnDevice()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve no space left on device: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

const (
	cmdPrefix = "$"
	dirPrefix = "dir"

	cdCommand = "cd"
	lsCommand = "ls"

	sizeThreshold = 100000
)

type terminalOutputParser struct {
	sc *bufio.Scanner

	lineNum  int
	nextLine []byte

	cmd *command
	err error
}

func newTerminalOutputParser(r io.Reader) *terminalOutputParser {
	return &terminalOutputParser{
		sc: bufio.NewScanner(r),
	}
}

func (p *terminalOutputParser) Next() bool {
	if p.err != nil {
		return false
	}

	if len(p.nextLine) == 0 {
		p.nextLine, p.err = p.readLine()
		if p.err != nil {
			return false
		}
	}

	oldLineLen := len(p.nextLine)
	p.nextLine = bytes.TrimPrefix(p.nextLine, []byte(cmdPrefix))
	if len(p.nextLine) == oldLineLen {
		p.err = fmt.Errorf("malformed command on line %d, must start with %q sign", p.lineNum, cmdPrefix)
		return false
	}

	elements := strings.Fields(string(p.nextLine))
	if len(elements) == 0 {
		p.err = fmt.Errorf("malformed command on line %d, empty command", p.lineNum)
		return false
	}

	p.cmd = &command{
		Name: elements[0],
		Args: elements[1:],
	}

	for {
		p.nextLine, p.err = p.readLine()
		if p.err != nil {
			return false
		}

		if bytes.HasPrefix(p.nextLine, []byte(cmdPrefix)) {
			return true
		}

		p.cmd.Output = append(p.cmd.Output, p.nextLine)
	}
}

func (p *terminalOutputParser) readLine() ([]byte, error) {
	p.lineNum++
	if p.sc.Scan() {
		return p.sc.Bytes(), nil
	}

	if err := p.sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error on line %d: %v", p.lineNum, err)
	}
	return nil, io.EOF
}

func (p *terminalOutputParser) Command() *command {
	if p.err != nil {
		panic("called Command while Next returned false")
	}
	return p.cmd
}

func (p *terminalOutputParser) Err() error {
	if p.err == io.EOF {
		return nil
	}
	return p.err
}

type command struct {
	Name   string
	Args   []string
	Output [][]byte
}

// fsDir is a dir-only version of filesystem.
type fsDir struct {
	Name     string
	Size     int
	Visited  bool   // protection against double visiting.
	Prev     *fsDir // .. analogue.
	Children []*fsDir
}

func solveNoSpaceLeftOnDevice() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	parser := newTerminalOutputParser(f)

	root := &fsDir{
		Name: "",
	}
	var curDir *fsDir
	for parser.Next() {
		cmd := parser.Command()
		switch cmd.Name {
		case cdCommand:
			if len(cmd.Args) != 1 {
				return 0, fmt.Errorf("expected 1 argument for %q-command", cdCommand)
			}
			dir := cmd.Args[0]
			if strings.HasPrefix(dir, "/") && len(dir) != 1 {
				return 0, fmt.Errorf("absolute for %q-command isn't supported", cdCommand)
			}

			switch dir {
			case "/":
				curDir = root
			case "..":
				if curDir.Prev == nil {
					return 0, fmt.Errorf("can't %q to move out one level", cdCommand)
				}
				curDir = curDir.Prev
			default:
				found := false
				for _, child := range curDir.Children {
					if child.Name == dir {
						curDir = child
						found = true
						break
					}
				}
				if !found {
					return 0, fmt.Errorf("current wd doesn't have %s according to registry", dir)
				}
			}
		case lsCommand:
			if len(cmd.Args) != 0 {
				return 0, fmt.Errorf("args for %q-command aren't supported", lsCommand)
			}
			if curDir.Visited {
				// I assume that fs traversing visits each directory only once.
				return 0, fmt.Errorf("already visited %s", curDir.Name)
			}

			for _, line := range cmd.Output {
				chunks := bytes.SplitN(line, []byte(" "), 2)
				if len(chunks) < 2 {
					return 0, fmt.Errorf("malformed %q output", lsCommand)
				}
				switch string(chunks[0]) {
				case "dir":
					dir := string(chunks[1])
					for _, child := range curDir.Children {
						if child.Name == dir {
							// I assume that fs traversing visits each directory only once.
							return 0, fmt.Errorf("directory %s already exists in registry", dir)
						}
					}

					curDir.Children = append(curDir.Children, &fsDir{
						Name: dir,
						Prev: curDir,
					})
				default:
					size, err := strconv.Atoi(string(chunks[0]))
					if err != nil {
						return 0, fmt.Errorf("wrong file %s size: %v", chunks[1], err)
					}

					// Set size for current wd and all previous.
					cur := curDir
					for cur != nil {
						cur.Size += size
						cur = cur.Prev
					}
				}
			}

			curDir.Visited = true
		default:
			return 0, fmt.Errorf("unknown command %q", cmd.Name)
		}
	}
	if err = parser.Err(); err != nil {
		return 0, fmt.Errorf("parse term output: %v", err)
	}

	return calculateSizesWithThresholdSum(root), nil
}

func calculateSizesWithThresholdSum(root *fsDir) int {
	sum := 0
	if root.Size <= sizeThreshold {
		sum += root.Size
	}

	for _, child := range root.Children {
		sum += calculateSizesWithThresholdSum(child)
	}

	return sum
}
