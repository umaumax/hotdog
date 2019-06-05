package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

var (
	firstRegexpFlag   string
	middleRegexpFlag  string
	lastRegexpFlag    string
	lineSeparatorFlag string
	verboseFlag       bool
)

func init() {
	flag.StringVar(&firstRegexpFlag, "first", "", "first line regexp")
	flag.StringVar(&middleRegexpFlag, "middle", ".*", "middle line regexp")
	flag.StringVar(&lastRegexpFlag, "last", "", "last line regexp")
	flag.StringVar(&lineSeparatorFlag, "separator", string(0x1e), "line separator default is '0x1e'(Record Separator)")
	flag.BoolVar(&verboseFlag, "verbose", false, "verbose flag")
}

type FilterStatus int

const (
	FilterStatusFirst FilterStatus = iota
	FilterStatusMiddle
	FilterStatusLast
)

var FilterStatusName = []string{"first", "middle", "last"}

type Filter struct {
	FirstRegexp  *regexp.Regexp
	MiddleRegexp *regexp.Regexp
	LastRegexp   *regexp.Regexp
	nextState    FilterStatus
	Separator    string
	blocks       [][]string
	lines        []string
}

func (f *Filter) Parse(line string) {
	if verboseFlag {
		log.Printf("[VERBOSE]: parse line '%s'\n", line)
	}
	switch f.nextState {
	case FilterStatusFirst:
		ret := f.FirstRegexp.MatchString(line)
		if !ret {
			f.lines = nil
			return
		}
		f.nextState = FilterStatusMiddle
	case FilterStatusMiddle:
		// NOTE: last regexp is stronger than middle
		if ret := f.LastRegexp.MatchString(line); ret {
			f.nextState = FilterStatusLast
			f.Parse(line)
			return
		}
		ret := f.MiddleRegexp.MatchString(line)
		if !ret {
			f.nextState = FilterStatusLast
			f.Parse(line)
			return
		}
	case FilterStatusLast:
		ret := f.LastRegexp.MatchString(line)
		f.nextState = FilterStatusFirst
		if !ret {
			f.lines = nil
			return
		}
	}
	f.lines = append(f.lines, line)
	if verboseFlag {
		log.Printf("[VERBOSE]: match line '%s'\n", line)
		log.Printf("[VERBOSE]: next state is '%s'\n", FilterStatusName[f.nextState])
	}

	if f.nextState == FilterStatusFirst {
		f.blocks = append(f.blocks, f.lines)
	}
}

func (f *Filter) String() string {
	w := &bytes.Buffer{}
	for _, block := range f.blocks {
		for _, line := range block {
			fmt.Fprintf(w, "%s%s", line, f.Separator)
		}
		fmt.Fprintf(w, "\n")
	}
	return w.String()
}

func main() {
	flag.Parse()
	firstRegexp, err := regexp.Compile(firstRegexpFlag)
	if err != nil {
		log.Fatalf("--first regexp '%s' is invalid", firstRegexpFlag)
	}
	middleRegexp, err := regexp.Compile(middleRegexpFlag)
	if err != nil {
		log.Fatalf("--middle regexp '%s' is invalid", middleRegexpFlag)
	}
	lastRegexp, err := regexp.Compile(lastRegexpFlag)
	if err != nil {
		log.Fatalf("--last regexp '%s' is invalid", lastRegexpFlag)
	}

	filter := &Filter{
		FirstRegexp:  firstRegexp,
		MiddleRegexp: middleRegexp,
		LastRegexp:   lastRegexp,
		Separator:    lineSeparatorFlag,
	}
	// NOTE: default input file is input pipe
	var inputFiles []string
	if flag.NArg() == 0 {
		inputFiles = append(inputFiles, os.Stdin.Name())
	} else {
		inputFiles = append(inputFiles, flag.Args()...)
	}
	for _, inputFile := range inputFiles {
		if verboseFlag {
			log.Printf("[VERBOSE]: input file '%s'\n", inputFile)
		}
		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		n := 0
		for scanner.Scan() {
			line := scanner.Text()
			filter.Parse(line)
			n++
		}
		if err = scanner.Err(); err != nil {
			return
		}
	}
	fmt.Print(filter.String())
}
