package internal

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ergagnon/ginder/infrastructure"
)

type app struct {
	Config AppConfig
}

func NewApp(cfg AppConfig) *app {
	return &app{
		Config: cfg,
	}
}

func (a *app) Run() {
	start := time.Now()

	fmt.Printf(a.Config.Directory)

	service := infrastructure.NewRawTextService()
	defer service.Close()

	extractPool := sync.Pool{
		New: func() any {
			return &extract{service: service}
		},
	}

	var wg sync.WaitGroup

	err := filepath.Walk(a.Config.Directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, ".ini") {
			return nil
		}

		extract := extractPool.Get().(*extract)

		wg.Add(1)
		go func() {
			extract.Extract(path)
			wg.Done()
		}()

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory %q: %v\n", a.Config.Directory, err)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Extract took %s", elapsed)
}

type extract struct {
	service infrastructure.RawTextService
}

func (me *extract) Extract(filePath string) {
	log.SetOutput(os.Stdout)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error open file %s", filePath)
		return
	}
	defer file.Close()

	log.Println("Extracting file: ", filePath)

	reader := me.service.Extract(file)

    // Create a pipe for async processing
	pr, pw := io.Pipe()

	var wg sync.WaitGroup

	// Goroutine to copy from original reader to pipe writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pw.Close()

		_, err := io.Copy(pw, reader)
		if err != nil {
			log.Printf("error reading data: %v", err)
			pw.CloseWithError(err)
		}
	}()

    // Goroutine to copy from pipe reader to file
	wg.Add(1)
	go func() {
		defer wg.Done()

		outputFile, err := os.Create(filePath + ".extracted")
		if err != nil {
			log.Printf("error creating output file: %v", err)
			return
		}
		defer outputFile.Close()

		_, err = io.Copy(outputFile, pr)
		if err != nil {
			log.Printf("error writing to file: %v", err)
		}
	}()

	wg.Wait()
	log.Println("Extraction completed successfully")
}
