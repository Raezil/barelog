
# barelog

**barelog** — это минималистичная библиотека логирования для Go без внешних зависимостей. 
Поддерживает цветной вывод, уровни логирования, глобальный логгер и логгер в контексте.

> 

---

## Установка

```bash
go get github.com/buraev/barelog@latest
```

---

## Быстрый старт

```go
package main

import "github.com/buraev/barelog"

func main() {
    barelog.Init() // автоматически настраивает логгер из переменных окружения

    barelog.Info("Сервер запущен", "port", 8080)
    barelog.Debug("Подробности", "trace_id")
}
```

---

## Уровни логирования

* `DEBUG`
* `INFO`
* `WARN`
* `ERROR`

Выводятся в цвете (если терминал поддерживает ANSI):

```go
barelog.Debug("отладочная информация")
barelog.Info("инфо-сообщение", "user", "alice")
barelog.Warn("предупреждение", "attempt", 2)
barelog.Error("ошибка", "err", "connection refused")
```

---

## Глобальный логгер

Вы можете использовать barelog напрямую без создания `Logger` вручную:

```go
barelog.Info("готов к работе")
```

Если нужно — установите свой глобальный логгер:

```go
logger := barelog.New(barelog.DEBUG)
barelog.SetGlobal(logger)
```

---

## Логгер в контексте

Подходит для middleware, request-scope логирования и т. д.:

```go
ctx := barelog.WithContext(context.Background(), barelog.New(barelog.DEBUG))
log := barelog.FromContext(ctx)
log.Info("из контекста", "req_id")
```

---

## Настройка через переменные окружения

| Переменная      | Описание            | Пример          |
| --------------- | ------------------- | --------------- |
| `BARELOG_LEVEL` | Уровень логирования | `debug`, `info` |

Пример:

```bash
BARELOG_LEVEL=debug ./yourApp
```

---

## Возможности ( планируются )

*

---

## Почему `barelog`?

* Zero dependencies
* Один файл — легко читать, понимать и адаптировать
* Хорош для CLI, сервисов, микросервисов, тестов

---

## Лицензия

MIT

---

## Автор

[github.com/buraev](https://github.com/buraev)
