package parser

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"text/template"
)

type Parser struct {
	regex *regexp.Regexp
	tmpl  *template.Template
}

type FuncInfo struct {
	FuncName   string
	FuncParams []Field
	FuncReturn []Field
}

type Field struct {
	Name     string
	Type     string
	Optional bool
}

//func (i *Info) SyncPermitMap() error
func NewParser() *Parser {
	reg, err := regexp.Compile("func\\s+(\\(\\w+\\s+\\*?\\w+\\)\\s+)*(\\w+)\\((.*?)\\)(\\s+.+)*")
	if err != nil {
		panic("regexp compile error:" + err.Error())
	}
	t := template.Must(template.New("commgen").Parse(CommentTemplate))

	return &Parser{
		regex: reg,
		tmpl:  t,
	}
}

func (p *Parser) ParseLine(line string) FuncInfo {
	fi := FuncInfo{}
	result := p.regex.FindStringSubmatch(line)
	for i, v := range result {
		if i == 2 {
			fi.FuncName = v
		}
		if i == 3 {
			fi.FuncParams = p.ParseParams(v)
		}
		if i == 4 {
			fi.FuncReturn = p.ParseReturn(v)
		}
	}
	return fi
}

func (p *Parser) ParseParams(line string) []Field {
	fieldInfo := strings.Split(line, ", ")
	result := make([]Field, 0, len(fieldInfo))

	tempName := []string{}
	for _, v := range fieldInfo {
		temp := strings.Split(v, " ")
		if len(temp) == 1 {
			tempName = append(tempName, temp[0])
		}
		if len(temp) == 2 {
			if len(tempName) != 0 {
				for _, tn := range tempName {
					result = append(result, Field{Name: tn, Type: temp[1]})
				}
				tempName = tempName[:0]
			}
			result = append(result, Field{Name: temp[0], Type: temp[1]})
		}
	}
	return result
}

func (p *Parser) ParseReturn(line string) []Field {
	l := strings.TrimSpace(line)
	if strings.HasPrefix(l, "(") && strings.HasSuffix(l, ")") {
		l = l[1 : len(l)-1]
	}
	result := make([]Field, 0, 0)
	f := strings.Split(l, ", ")
	for _, v := range f {
		temp := strings.Split(v, " ")
		if len(temp) == 2 {
			result = append(result, Field{Name: temp[0], Type: temp[1]})
		}
		if len(temp) == 1 && temp[0] != "" {
			result = append(result, Field{Type: temp[0]})
		}
	}
	return result
}

func (p *Parser) Exec(fi FuncInfo) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := p.tmpl.Execute(buf, fi)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Parser) Parse(line string) ([]byte, error) {
	if strings.HasPrefix(line, "func") {
		if strings.HasSuffix(line, "{") {
			line = line[0 : len(line)-1]
		}

		fi := p.ParseLine(line)
		if fi.FuncName != "" {
			return p.Exec(p.ParseLine(line))
		} else {
			return nil, errors.New("invalid line")
		}
	} else {
		return nil, nil
	}
}
