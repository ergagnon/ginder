package internal

import (
	"fmt"
	"io"
	"log"
	"os"
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
    var wg sync.WaitGroup

    start := time.Now()
	fmt.Printf(a.Config.Message)

    service := infrastructure.NewRawTextService()
    defer service.Close()

    filePath := "C:\\Users\\egagn\\OneDrive\\Documents\\Inspiration.docx"

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("error open file %s", filePath)
        return
    }
    defer file.Close()

    reader := service.Extract(file)

    buf := make([]byte, 1024)

    wg.Add(1)
    go func() {
        for {
            n, err := reader.Read(buf)

            if err == io.EOF {
                wg.Done()
                break
            }

            log.Println(n, err, string(buf[:n]))
        }
    }()

    wg.Wait()

    elapsed := time.Since(start)
    log.Printf("Extract took %s", elapsed)
}
