package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"go.uber.org/zap"

    "github.com/yukkat/practicum-alice-skill/internal/logger"
    "github.com/yukkat/practicum-alice-skill/internal/models"
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
    return http.ListenAndServe(flagRunAddr, logger.RequestLogger(GzipMiddleware(webhook)))
}

func webhook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // десериализуем запрос в структуру модели
    logger.Log.Debug("decoding request")
    var req models.Request
    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&req); err != nil {
        logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // проверяем, что пришёл запрос понятного типа
    if req.Request.Type != models.TypeSimpleUtterance {
        logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    // заполняем модель ответа
    resp := models.Response{
        Response: models.ResponsePayload{
            Text: "Извините, я пока ничего не умею",
        },
        Version: "1.0",
    }

    w.Header().Set("Content-Type", "application/json")

    // сериализуем ответ сервера
    enc := json.NewEncoder(w)
    if err := enc.Encode(resp); err != nil {
        logger.Log.Debug("error encoding response", zap.Error(err))
        return
    }
    logger.Log.Debug("sending HTTP 200 response")
}