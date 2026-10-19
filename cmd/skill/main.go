package main

import (
	"fmt"
	"net/http"

	"github.com/yukkat/practicum-alice-skill/internal/logger"
	"go.uber.org/zap"
)

func main() {
    // обрабатываем аргументы командной строки
    parseFlags()

    if err := run(); err != nil {
        panic(err)
    }
}

func run() error {
    if err := logger.Initialize(flagLogLevel); err != nil {
        return err
    }

    logger.Log.Info("Running server", zap.String("address", flagRunAddr))
    fmt.Println("Running server on", flagRunAddr)
    return http.ListenAndServe(flagRunAddr, logger.RequestLogger(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`
          {
            "response": {
              "text": "Извините, я пока ничего не умею"
            },
            "version": "1.0"
          }
        `))
    logger.Log.Debug("sending HTTP 200 response")
}