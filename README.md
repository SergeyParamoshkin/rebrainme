# Context

## Что такое context в Go?

context.Context — это стандартный интерфейс в Go, предназначенный для:

отмены операций
установки таймаутов и дедлайнов
передачи метаданных по цепочке вызовов (например, в middleware)

## Зачем он нужен?

Управление временем выполнения операций (например, прерывание запроса к БД).
Уменьшение утечек горутин.
Общий механизм отмены операций в сетевых вызовах, асинхронных задачах.

## Где используется?

В стандартной библиотеке (http, database/sql, net/http, os/exec)
В gRPC
Везде

# Основные концепции

```
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

* Done() — канал, который закрывается при отмене или истечении времени.
* Err() — возвращает причину отмены (context.Canceled, context.DeadlineExceeded).
* Deadline() — возвращает дедлайн, если установлен.
* Value() — возвращает значение, связанное с ключом

# Создание контекста

* context.Background()
Базовый пустой контекст — стартовая точка.
Используется в main, и в тестах

* context.TODO()
Заполнитель, когда контекст нужен, но неизвестно, какой.

* context.WithCancel(parent)
Создаёт контекст, который можно отменить вручную:

```
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

```
context.WithTimeout(parent, duration)
```

Контекст отменяется по истечении таймаута:

```
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

```
context.WithDeadline(parent, time)
```

Аналогично, но на конкретное время.

```
context.WithValue(parent, key, value)
```

Контекст с дополнительным значением:

```
ctx := context.WithValue(ctx, "requestID", "abc123")
```

# Примеры

```
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    select {
    case <-time.After(2 * time.Second):
        fmt.Fprintln(w, "done")
    case <-ctx.Done():
        http.Error(w, "request cancelled", http.StatusRequestTimeout)
    }
}
```

```
func fetchUser(ctx context.Context, db *sql.DB, id int) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    row := db.QueryRowContext(ctx, "SELECT name FROM users WHERE id=?", id)
    ...
}
```

### До 1.20

Ранее, до Go 1.20, чтобы понять, почему контекст был отменён, приходилось использовать:
`ctx.Err()`
Но ctx.Err(), всегда возвращает только одно из двух значений (Canceled, DeadlineExceeded),
не позволяет узнать вложенные ошибки, если они были причиной отмены.

Что делает context.Cause()?

`context.Cause(ctx)` возвращает реальную причину отмены контекста, включая:

ошибки, переданные в WithCancelCause
context.Canceled, context.DeadlineExceeded
любую другую ошибку, которую мы сами передали

```
ctx, cancel := context.WithCancelCause(context.Background())
cancel(errors.New("user aborted operation"))

<-ctx.Done()

fmt.Println(context.Cause(ctx)) // "user aborted operation"
fmt.Println(ctx.Err())          // "context canceled"

```

```
func handler(ctx context.Context) error {
    ctx, cancel := context.WithCancelCause(ctx)
    defer cancel(errors.New("request aborted by handler"))

    // что-то делаем...

    <-ctx.Done()

    return context.Cause(ctx)
}

```

# Плохие практики

```
// Плохо:
ctx := context.WithValue(context.Background(), "userID", 42)
// Хорошо:
func handle(ctx context.Context, userID int) {}
```

```
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// cancel забыли -> ресурсы не освободятся
```

```
select {
case <-ctx.Done():
    // Плохо: не обрабатывается err
}
```

# И ещё

```
type Handler struct {
    ctx context.Context // плохо, Контекст должен передаваться явно в аргументах.
}
```

Не передавай nil вместо контекста
— Всегда передавай context.Background() хотя бы.

Не злоупотребляй context.WithValue
— Только для request-scoped данных: requestID, auth token.

Никогда не забывай cancel()
— Особенно при использовании WithTimeout, WithCancel.

Не используй context.TODO() в production-коде, старайтесь)
