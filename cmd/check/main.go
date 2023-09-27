package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/hashicorp/go-multierror"

	"procinspect/pkg/checker"
	"procinspect/pkg/log"
	"procinspect/pkg/semantic"
)

var (
	file = flag.String("file", "", "")
	dir  = flag.String("dir", "", "")
	prof = flag.Bool("prof", false, "")
	bin  = flag.String("bin", "", "")
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
	{
		i := 0
		if *dir != "" {
			i++
		}
		if *file != "" {
			i++
		}
		if *bin != "" {
			i++
		}
		if i > 1 {
			fmt.Println("only one of -dir, -file, -bin is allowed")
			return
		}
	}

	var sql string
	if *dir != "" {
		// walk directory
		err := filepath.Walk(*dir, walkDir)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("parse ok")
			return
		}

	} else if *file != "" {
		// get abstract path of file
		sql = *file
		_ = parseFile(sql)
		return
	} else if *bin != "" {
		filename := *bin
		checkBinary(filename)
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

func parseFile(sql string) error {
	absPath, err := filepath.Abs(sql)
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
	// parse file
	fmt.Print("parse ", filepath.Base(absPath), " ")
	start := time.Now()
	script, err := checker.LoadScript(string(text))
	elapsed := time.Since(start)
	if err != nil {
		// name := filepath.Base(absPath)
		fmt.Printf("error:\n%s\n", err)
		return err
	} else {
		fmt.Print("ok;")
		fmt.Printf(" time: %s\n", elapsed)
	}

	return check(script)
}

func check(script *semantic.Script) error {
	v := checker.NewValidVisitor()
	_ = script.Accept(v)

	var errs *multierror.Error
	errors.As(v.Error(), &errs)
	for _, e := range errs.Errors {
		err := e.(checker.SqlValidationError)
		log.Warn("unsupported", log.String("err", err.Error()), log.Int("line", err.Line))
	}
	return nil
}

func checkBinary(filename string) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open(absPath)
	if err != nil {
		return
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return
	}
	defer gr.Close()

	// read file to string
	text, err := io.ReadAll(gr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// parse file
	log.Info("Start Parse", log.String("file", filepath.Base(absPath)))
	start := time.Now()
	script, err := semantic.NewNodeDecoder[*semantic.Script]().Decode(text)
	elapsed := time.Since(start)
	log.Info("End Parse", log.String("file", filepath.Base(absPath)),
		log.String("duration", elapsed.String()))
	if err != nil {
		// name := filepath.Base(absPath)
		fmt.Printf("error:\n%s\n", err)
		return
	}
	check(script)
}
