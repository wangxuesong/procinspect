package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"procinspect/pkg/log"
	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

var (
	file     = flag.String("file", "", "")
	dir      = flag.String("dir", "", "")
	prof     = flag.Bool("prof", false, "")
	parallel = flag.Bool("p", false, "parallel parse")
	verbose  = flag.Bool("v", false, "verbose")
	size     = flag.Int64("size", 500*1024, "size")
)

type (
	ParseRequest struct {
		FileName string
		Source   string
		Index    int
		Start    int
	}
	ParseResult struct {
		FileName string
		Index    int
		Start    int
		AstFunc  func(int) (*semantic.Script, error)
		Error    error
		Source   string
	}
)

func main() {
	flag.Parse()
	//flag.PrintDefaults()

	if *prof {
		pf, err := os.Create("./cpu.prof")
		if err != nil {
			fmt.Printf("创建采集文件失败, err:%v\n", err)
			return
		}
		pprof.StartCPUProfile(pf)
		defer pprof.StopCPUProfile()
	}

	var sql string
	if *dir != "" {
		log.Info("Start Check", log.String("dir", *dir))
		// walk directory
		err := filepath.Walk(*dir, walkDir)
		if err != nil {
			fmt.Println(err)
			log.Error("Check Error", log.String("error", err.Error()))
			return
		} else {
			log.Info("End Check", log.String("dir", *dir))
			return
		}

	} else if *file != "" {
		// get abstract path of file
		sql = *file
		_ = parseFile(sql)
		return
	}
}

func walkDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	return parseFile(path)
}

func parseFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// read file to string
	text, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	count := make(chan int)
	go func() {
		lines := strings.Split(string(text), "\n")
		count <- len(lines)
	}()
	// parse file
	log.Info("Start Parse", log.String("file", filepath.Base(absPath)))

	requests, err := prepareRequest(absPath)

	var results []*ParseResult = make([]*ParseResult, len(requests))
	start := time.Now()
	if *parallel {
		parseChan := make(chan *ParseResult)
		wg := &sync.WaitGroup{}
		// runtime.GOMAXPROCS(0)
		for _, req := range requests {
			wg.Add(1)
			go func(req *ParseRequest) {
				defer wg.Done()
				parseChan <- parseBlock(req)
			}(req)
		}
		for _, _ = range requests {
			result := <-parseChan
			results[result.Index] = result
		}
		close(parseChan)
		wg.Wait()
	} else {
		for _, req := range requests {
			result := parseBlock(req)
			results[result.Index] = result
			//results = append(results, result)
		}
	}
	elapsed := time.Since(start)
	lines := <-count
	for _, result := range results {
		var err error
		if result.Error != nil {
			log.Error("Parse Error", log.String("file", filepath.Base(absPath)),
				log.String("duration", elapsed.String()),
				log.Int("index", result.Index),
				log.Int("start", result.Start),
				log.Int("lines", lines),
				log.String("error", result.Error.Error()),
				log.String("source", result.Source),
			)
			err = errors.Join(err, result.Error)
		}
		if err != nil {
			return err
		}
	}

	log.Info("End Parse", log.String("file", filepath.Base(absPath)),
		log.Int("lines", lines),
		log.String("duration", elapsed.String()))

	// 生成 AST
	for _, result := range results {
		_, err := result.AstFunc(result.Start)
		if err != nil {
			log.Error("Check Error", log.String("file", filepath.Base(absPath)),
				log.String("error", err.Error()),
			)
		}
	}
	return nil
}

func prepareRequest(path string) ([]*ParseRequest, error) {
	requests := make([]*ParseRequest, 0)

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	name := filepath.Base(path)
	text, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	source := string(text)

	if *parallel && info.Size() > *size {
		re := regexp.MustCompile(`\r\n`)
		source = re.ReplaceAllString(source, "\n")
		// split source by /
		regex := regexp.MustCompile(`(?m)^/$`)
		blocks := regex.Split(source, -1)
		start := 0
		offset := 0
		for i, block := range blocks {
			requests = append(requests, &ParseRequest{
				FileName: name,
				Source:   block,
				Index:    i,
				Start:    start + offset,
			})
			start += strings.Count(block, "\n")
			offset = 0
		}
	} else {
		requests = append(requests, &ParseRequest{
			FileName: name,
			Source:   source,
			Index:    0,
			Start:    0,
		})
	}

	return requests, nil
}

func parseBlock(r *ParseRequest) *ParseResult {
	result := &ParseResult{
		FileName: r.FileName,
		Index:    r.Index,
		Start:    r.Start,
	}
	if *verbose {
		result.Source = r.Source
	}
	s, err := parser.ParseSql(r.Source)
	if err != nil {
		result.Error = err
	} else {
		result.AstFunc = s
	}
	return result
}
