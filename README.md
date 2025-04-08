
# slog: Gelf (Graylog) handler

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-%23007d9c)
[![License](https://img.shields.io/github/license/samber/slog-graylog)](./LICENSE)

A [Graylog](https://www.graylog.org/) Handler for [slog](https://pkg.go.dev/log/slog) Go library.

## 🚀 Install

```sh
go get github.com/eandr-67/gelf_slog
```

**Compatibility**: go >= 1.24

## Причины создания

Исходный вариант пакета [] намертво прибит гвоздями к первой версии
пакета [Graylog2/go-gelf](https://github.com/Graylog2/go-gelf).
Но существует вторая версия - тоже древняя но всё же более функциональная.

Хотелось и перейти на вторую версию транспортного пакета, и
лучше приспособить адаптер для DI-контейнеров.
Кроме того, в текущей реализации есть проблема 

К сожалению, полностью отвязать адаптер от конкретной реализации транспорта 

Для этого пришлось перейти от типа к интерфейсу.
В сам же интерфейс добавить    
и отвязать адаптер от транспорта. Для чего надо перейти от типа
к интерфейсу - в обоих пакетах. Так что пришлось делать два форка:
этот и [enadr-67/gelf](https://github.com/eandr-67/gelf). 



## 💡 Usage

### Handler options

```go
type Option struct {
    // log level (default: debug)
    Level slog.Leveler

    // connection to graylog
    Writer gelf.Writer

    // optional: customize json payload builder
    Converter Converter
    // optional: fetch attributes from context
    AttrFromContext []func(ctx context.Context) []slog.Attr

    // optional: see slog.HandlerOptions
    AddSource   bool
    ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}
```

Attributes will be injected in log "extra".

Other global parameters:

```go
gelf_slog.SourceKey = "source"
gelf_slog.ErrorKeys = []string{"error", "err"}
gelf_slog.LogLevels = map[slog.Level]int32{...}
```

### Example

```go
import (
    "github.com/eandr-67/gelf"
    "github.com/eandr-67/gelf_log"
    "log/slog"
)

func main() {
    // docker-compose up -d
    // or
    // ncat -l 12201 -u
    gelfWriter, err := gelf.NewUDPWriter("localhost:12201")
    if err != nil {
        log.Fatalf("gelf.NewWriter: %s", err)
    }

	gelfWriter.CompressionType = gelf.CompressNone  // for debugging only

    logger := slog.New(sloggraylog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler())
    logger = logger.
        With("environment", "dev").
        With("release", "v1.0.0")

    // log error
    logger.
        With("category", "sql").
        With("query.statement", "SELECT COUNT(*) FROM users;").
        With("query.duration", 1*time.Second).
        With("error", fmt.Errorf("could not count users")).
        Error("caramba!")

    // log user signup
    logger.
        With(
            slog.Group("user",
                slog.String("id", "user-123"),
                slog.Time("created_at", time.Now()),
            ),
        ).
        Info("user registration")
}
```

Output:

```json
{
    "timestamp":"2023-04-10T14:00:0.000000+00:00",
    "level":3,
    "message":"caramba!",
    "extra":{
        "error":{
            "error":"could not count users",
            "kind":"*errors.errorString",
            "stack":null
        },
        "environment":"dev",
        "release":"v1.0.0",
        "category":"sql",
        "query.statement":"SELECT COUNT(*) FROM users;",
        "query.duration": "1s"
    }
}


{
    "timestamp":"2023-04-10T14:00:0.000000+00:00",
    "level":6,
    "message":"user registration",
    "extra":{
        "environment":"dev",
        "release":"v1.0.0",
        "user":{
            "id":"user-123",
            "created_at":"2023-04-10T14:00:0.000000+00:00"
        }
    }
}
```

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
