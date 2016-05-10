package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Packages map[string]string `yaml:"packages"`
}

type ImportInfo struct {
	StartLine int
	EndLine   int
}

const (
	GENERATE_IMPORT_START_COMMENT = "+imports"
	GENERATE_IMPORT_END_COMMENT   = "+imports-end"
)

var (
	configFile = flag.String("c", "import.yml", "set config file")
)

func main() {
	flag.Parse()

	args := flag.Args()

	file := ""
	if len(args) == 0 {
		file = "./main.go"
	}

	config, err := NewConfig(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	info, err := Parse(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = PkgImport(file, config, info)
}

func NewConfig(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func Parse(file string) (*ImportInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	ast.FileExports(f)

	var start, end int
	for _, comment := range f.Comments {
		switch strings.Trim(comment.Text(), "\n") {
		case GENERATE_IMPORT_START_COMMENT:
			file := fset.File(comment.Pos())
			start = file.Line(comment.Pos())

		case GENERATE_IMPORT_END_COMMENT:
			file := fset.File(comment.Pos())
			end = file.Line(comment.Pos())
		default:
		}

	}

	info := &ImportInfo{
		StartLine: start,
		EndLine:   end,
	}

	return info, nil
}

func PkgImport(file string, conf *Config, info *ImportInfo) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var inputLines []string
	inputCount := 0

	scanner := bufio.NewScanner(bytes.NewReader(input))

	// remove already generate import line
	for scanner.Scan() {
		inputCount++

		if info.StartLine < inputCount && inputCount < info.EndLine {
			continue
		}

		inputLines = append(inputLines, scanner.Text())
	}

	// insert config packages
	var outputLines []string
	for i, v := range inputLines {
		outputLines = append(outputLines, v)

		if i+1 == info.StartLine {
			outputLines = append(outputLines, "\n")

			for k, v := range conf.Packages {
				pkg := ""
				if k == "" {
					pkg = fmt.Sprintf("\"%s\"", v)
				} else {
					pkg = fmt.Sprintf("%s \"%s\"", k, v)
				}

				outputLines = append(outputLines, pkg)
			}
		}
	}

	f, err := os.Create("./main.go")
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for _, line := range outputLines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
