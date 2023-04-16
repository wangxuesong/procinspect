package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"procinspect/pkg/parser"
)

var (
	file = flag.String("file", "", "")
	dir  = flag.String("dir", "", "")
	prof = flag.Bool("prof", false, "")
)

func main() {
	flag.Parse()
	flag.PrintDefaults()

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
	err = parser.Parse(string(text))
	elapsed := time.Since(start)
	if err != nil {
		//name := filepath.Base(absPath)
		fmt.Printf("error: %s\n", err)
		return err
	} else {
		fmt.Print("ok;")
		fmt.Printf(" time: %s\n", elapsed)
		return nil
	}
}
