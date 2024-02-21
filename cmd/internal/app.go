package internal

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	var wg sync.WaitGroup

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error open file %s", filePath)
		return
	}
	defer file.Close()

	log.Println("Extracting file: ", filePath)

	reader := me.service.Extract(file)

	buf := make([]byte, 1024)

	wg.Add(1)
	go func() {
		for {
			_, err := reader.Read(buf)

			if err == io.EOF {
				wg.Done()
				break
			}

			//log.Println(n, err, string(buf[:n]))
		}
	}()

	wg.Wait()
}
