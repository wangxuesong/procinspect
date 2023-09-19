package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"

	"atomicgo.dev/cursor"
	"github.com/pterm/pterm"
	"go.uber.org/zap"

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
	index    = flag.Int("index", 0, "start from index")
)

type WorkerPool struct {
	workers    []*Worker
	taskQueue  chan Task
	numWorkers int
	wg         sync.WaitGroup
}

type Task func()

type Worker struct {
	id      int
	workers *WorkerPool
}

func NewWorkerPool(numWorkers int, taskQueueSize int) *WorkerPool {
	pool := &WorkerPool{
		workers:    make([]*Worker, numWorkers),
		taskQueue:  make(chan Task, taskQueueSize),
		numWorkers: numWorkers,
	}

	for i := 0; i < numWorkers; i++ {
		worker := &Worker{
			id:      i + 1,
			workers: pool,
		}
		pool.workers[i] = worker
		go worker.start()
	}

	return pool
}

func (p *WorkerPool) Submit(task Task) {
	p.taskQueue <- task
}

func (w *Worker) start() {
	for task := range w.workers.taskQueue {
		w.workers.wg.Add(1)
		func() {
			defer func() {
				w.workers.wg.Done()
				// runtime.GC()
			}()
			task()
		}()
	}
}

type (
	msg struct {
		Msg    string
		Fields []zap.Field
	}
	ParseRequest struct {
		FileName string
		Source   string
		Index    int
		Start    int
		Msg      chan msg
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
	// flag.PrintDefaults()

	if *prof {
		pf, err := os.Create("./cpu.prof")
		if err != nil {
			fmt.Printf("创建采集文件失败, err:%v\n", err)
			return
		}
		pprof.StartCPUProfile(pf)
		defer pprof.StopCPUProfile()
	}

	processSignal()

	if *verbose {
		log.SetLevel(log.DebugLevel)
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

func processSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		cursor.Show()
		fmt.Println("Exit")
		os.Exit(-1)
	}()
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
	msgChan := make(chan msg)
	defer close(msgChan)
	start := time.Now()
	if *parallel {
		parallelParse(requests, msgChan, results)
	} else {
		for _, req := range requests[*index:] {
			result := parseBlock(req)
			results[result.Index] = result
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

func parallelParse(requests []*ParseRequest, msgChan chan msg, results []*ParseResult) {
	p, _ := pterm.DefaultProgressbar.
		WithTotal(len(requests)).
		WithMaxWidth(-1).
		WithTitle("Parse file").
		Start()

	numWorkers := runtime.GOMAXPROCS(0) - 1
	pool := NewWorkerPool(numWorkers, len(requests))
	parseChan := make(chan *ParseResult)
	for _, req := range requests[*index:] {
		tmpReq := *req
		tmpReq.Msg = msgChan
		pool.Submit(func() {
			parseChan <- parseBlock(&tmpReq)
		})
	}
	for _, _ = range requests[*index:] {
		result := <-parseChan
		results[result.Index] = result
		if result.Error != nil {
			log.Error("Parse Error",
				log.Int("index", result.Index),
				log.Int("start", result.Start),
				log.String("error", result.Error.Error()),
			)
		}
		p.Increment()
		//	result.AstFunc(result.Start)
	}
	close(parseChan)
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
			if strings.TrimSpace(block) == "" {
				continue
			}
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
	log.Debug("start parse", log.String("foo", "sss"), log.Int("index", r.Index), log.Int("start", r.Start))
	el := time.Now()
	s, err := parser.ParseSql(r.Source)
	du := time.Since(el)
	log.Debug("stop parse", log.String("foo", "xxx"), log.Int("index", r.Index), log.String("duration", du.String()), log.Int("start", r.Start))
	if err != nil {
		result.Error = err
	} else {
		result.AstFunc = s
	}
	return result
}
