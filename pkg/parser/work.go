package parser

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

func Work(file string) error {
	if !strings.HasSuffix(file, ".go") {
		return errors.New("not support file type")
	}

	if err := exec.Command("go", "fmt", file).Run(); err != nil {
		return err
	}

	p := NewParser()
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	output := input
	lines := strings.Split(string(input), "\n")
	lastLine := ""
	for _, line := range lines {
		if strings.HasPrefix(lastLine, "//") || strings.HasPrefix(lastLine, "*/") {
			lastLine = line
			continue
		}
		line = strings.TrimSpace(line)
		newline, err := p.Parse(line)
		if err == nil {
			newline = append(newline, []byte(line)...)
			//lines[i] = string(newline)
		}
		output = bytes.Replace(output, []byte(line), newline, -1)
		lastLine = line
	}

	//output := strings.Join(lines, "\n")
	//log.Println(output)
	err = ioutil.WriteFile(file, output, 0644)
	if err != nil {
		return err
	}

	if err := exec.Command("go", "fmt", file).Run(); err != nil {
		return err
	}

	return nil
}
