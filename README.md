# Тестирование в Go

Этот репозиторий содержит примеры и описание различных подходов к тестированию в Go.

## Быстрый старт

```bash
# Клонирование репозитория
git clone https://github.com/SergeyParamoshkin/rebrainme.git
cd rebrainme

# Установка зависимостей
go mod download

# Запуск всех тестов
go test ./...

# Запуск тестов с подробным выводом
go test -v ./...

# Запуск бенчмарков
go test -bench=. ./examples/03_benchmarks/

# Запуск фаззинга (30 секунд)
go test -fuzz=FuzzReverse -fuzztime=30s ./examples/05_fuzzing/
```

## Структура примеров

```
examples/
├── 01_basic/           - Базовые примеры тестов (calculator)
├── 02_table_driven/    - Табличные тесты (strings)
├── 03_benchmarks/      - Примеры бенчмарков
├── 04_testing_main/    - Тестирование функции main
├── 05_fuzzing/         - Примеры фаззинг-тестов
├── 06_advanced/        - Продвинутые примеры (моки, testify)
└── 07_testdata/        - Использование testdata и golden files
```

## Содержание

- [Примеры кода](#примеры-кода)
- [Типы тестирования](#типы-тестирования)
- [Структура тестов](#структура-тестов)
- [Табличные тесты](#табличные-тесты)
- [Способы запуска тестов](#способы-запуска-тестов)
- [Бенчмарки](#бенчмарки)
  - [Что такое b.N?](#что-такое-bn)
- [Тестирование main](#тестирование-main)
- [Фаззинг](#фаззинг)
- [Продвинутые техники](#продвинутые-техники)
- [Директория testdata](#директория-testdata)

## Примеры кода

### 01_basic - Базовые тесты

**Расположение:** `examples/01_basic/`

Простые примеры unit-тестов для функций калькулятора.

```bash
# Запуск
go test ./examples/01_basic/

# С подробным выводом
go test -v ./examples/01_basic/
```

**Что внутри:**
- Простые unit-тесты
- Тесты с подтестами (subtests)
- Обработка ошибок в тестах

### 02_table_driven - Табличные тесты

**Расположение:** `examples/02_table_driven/`

Примеры табличных тестов для работы со строками.

```bash
# Запуск всех тестов
go test ./examples/02_table_driven/

# Запуск конкретного теста
go test -run TestReverse ./examples/02_table_driven/
```

**Что внутри:**
- Табличные тесты с разными сценариями
- Работа с Unicode и emoji
- Тестирование валидации
- Проверка инвариантов

### 03_benchmarks - Бенчмарки

**Расположение:** `examples/03_benchmarks/`

Сравнение производительности различных реализаций.

```bash
# Запуск всех бенчмарков
go test -bench=. ./examples/03_benchmarks/

# С информацией о памяти
go test -bench=. -benchmem ./examples/03_benchmarks/

# Конкретный бенчмарк
go test -bench=BenchmarkStringConcat ./examples/03_benchmarks/

# Сохранить результаты
go test -bench=. -benchmem ./examples/03_benchmarks/ > bench_results.txt
```

**Что внутри:**
- Сравнение конкатенации строк (+ vs Builder vs Join)
- Сравнение алгоритмов (рекурсивный vs итеративный Fibonacci)
- Параллельные бенчмарки
- Измерение аллокаций памяти

### 04_testing_main - Тестирование main

**Расположение:** `examples/04_testing_main/`

Различные подходы к тестированию функции main.

```bash
# Запуск тестов
go test ./examples/04_testing_main/

# Запуск приложения
go run ./examples/04_testing_main/main.go -name Alice -repeat 2
```

**Что внутри:**
- Тестирование с подменой stdout/stderr
- Тестирование флагов командной строки
- Example-тесты
- Изоляция логики от main

### 05_fuzzing - Фаззинг

**Расположение:** `examples/05_fuzzing/`

Примеры фаззинг-тестов для поиска багов.

```bash
# Запуск фаззинга (будет работать пока не найдет ошибку)
go test -fuzz=FuzzReverse ./examples/05_fuzzing/

# С ограничением по времени
go test -fuzz=FuzzReverse -fuzztime=30s ./examples/05_fuzzing/

# Запуск обычных тестов (включая регрессионные фазз-тесты)
go test ./examples/05_fuzzing/
```

**Что внутри:**
- Фаззинг строковых функций
- Фаззинг парсеров (JSON, URL, key-value)
- Проверка инвариантов
- Поиск паник и крэшей

### 06_advanced - Продвинутые техники

**Расположение:** `examples/06_advanced/`

Примеры с использованием моков и библиотеки testify.

```bash
# Запуск
go test ./examples/06_advanced/

# С подробным выводом
go test -v ./examples/06_advanced/

# Пропуск длительных интеграционных тестов
go test -short ./examples/06_advanced/
```

**Что внутри:**
- Использование testify (assert, require)
- Моки с помощью testify/mock
- Тестирование сервисного слоя
- Интеграционные тесты

### 07_testdata - Testdata и Golden Files

**Расположение:** `examples/07_testdata/`

Примеры использования директории testdata для хранения тестовых данных и golden files.

```bash
# Запуск тестов
go test -v ./examples/07_testdata/

# Обновление golden files
go test ./examples/07_testdata/ -update
```

**Что внутри:**
- Чтение тестовых данных из файлов
- Табличные тесты с множеством файлов
- Golden files для тестирования генерации текста
- Фикстуры для переиспользования данных
- Обновление golden files через флаг
- Структура testdata директории с README

**Структура testdata:**
```
testdata/
├── README.md                      # Документация
├── valid_config.json              # Тестовые данные
├── invalid_config.json
├── production_config.json
├── fixtures/                      # Фикстуры
│   └── users.json
└── golden/                        # Эталонные результаты
    ├── simple_report.txt
    ├── empty_report.txt
    └── quarterly_report.txt
```

## Типы тестирования

### 1. Unit-тесты (Модульное тестирование)

Тестирование отдельных функций и методов в изоляции.

```go
package calculator

func Add(a, b int) int {
    return a + b
}
```

```go
package calculator

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5

    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}
```

### 2. Integration-тесты (Интеграционное тестирование)

Тестирование взаимодействия между компонентами системы.

```go
func TestDatabaseUserRepository(t *testing.T) {
    db := setupTestDatabase(t)
    defer db.Close()

    repo := NewUserRepository(db)
    user := &User{Name: "John", Email: "john@example.com"}

    err := repo.Create(user)
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    found, err := repo.FindByEmail("john@example.com")
    if err != nil {
        t.Fatalf("Failed to find user: %v", err)
    }

    if found.Name != user.Name {
        t.Errorf("got name %s, want %s", found.Name, user.Name)
    }
}
```

### 3. End-to-End тесты

Тестирование полного пути пользовательского сценария.

```go
func TestUserRegistrationFlow(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()

    client := &http.Client{}

    // Регистрация
    resp, err := client.Post(
        server.URL+"/register",
        "application/json",
        strings.NewReader(`{"email":"test@example.com","password":"secret"}`),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusCreated)
    }
}
```

## Структура тестов

### Базовая структура

```go
func TestFunctionName(t *testing.T) {
    // Arrange (подготовка)
    input := "test"
    expected := "TEST"

    // Act (действие)
    result := ToUpper(input)

    // Assert (проверка)
    if result != expected {
        t.Errorf("ToUpper(%q) = %q; want %q", input, result, expected)
    }
}
```

### Subtests (подтесты)

```go
func TestMath(t *testing.T) {
    t.Run("addition", func(t *testing.T) {
        result := Add(2, 3)
        if result != 5 {
            t.Errorf("got %d, want 5", result)
        }
    })

    t.Run("subtraction", func(t *testing.T) {
        result := Subtract(5, 3)
        if result != 2 {
            t.Errorf("got %d, want 2", result)
        }
    })
}
```

### Helper функции

```go
func TestUserValidation(t *testing.T) {
    t.Helper()

    assertValid := func(t *testing.T, user User) {
        t.Helper()
        if err := user.Validate(); err != nil {
            t.Errorf("expected user to be valid, got error: %v", err)
        }
    }

    user := User{Name: "John", Email: "john@example.com"}
    assertValid(t, user)
}
```

## Табличные тесты

Табличные тесты (table-driven tests) - это идиоматический способ тестирования в Go, позволяющий запустить одну и ту же логику с разными входными данными.

### Простой пример

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a        int
        b        int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed", -2, 3, 1},
        {"zeros", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Расширенный пример с ошибками

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name      string
        dividend  float64
        divisor   float64
        want      float64
        wantError bool
    }{
        {
            name:      "normal division",
            dividend:  10,
            divisor:   2,
            want:      5,
            wantError: false,
        },
        {
            name:      "division by zero",
            dividend:  10,
            divisor:   0,
            want:      0,
            wantError: true,
        },
        {
            name:      "negative result",
            dividend:  -10,
            divisor:   2,
            want:      -5,
            wantError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.dividend, tt.divisor)

            if tt.wantError {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }

            if got != tt.want {
                t.Errorf("Divide(%f, %f) = %f; want %f",
                    tt.dividend, tt.divisor, got, tt.want)
            }
        })
    }
}
```

### Табличные тесты с setup/teardown

```go
func TestUserRepository(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name:    "valid user",
            user:    &User{Name: "John", Email: "john@example.com"},
            wantErr: false,
        },
        {
            name:    "invalid email",
            user:    &User{Name: "John", Email: "invalid"},
            wantErr: true,
        },
        {
            name:    "empty name",
            user:    &User{Name: "", Email: "test@example.com"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db := setupTestDB(t)
            defer db.Close()

            repo := NewUserRepository(db)
            err := repo.Create(tt.user)

            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Способы запуска тестов

### Основные команды

```bash
# Запуск всех тестов в текущей директории
go test

# Запуск всех тестов в проекте
go test ./...

# Verbose режим (подробный вывод)
go test -v

# Запуск конкретного теста
go test -run TestFunctionName

# Запуск тестов с определенным паттерном
go test -run TestUser

# Запуск конкретного подтеста
go test -run TestMath/addition

# Запуск с race detector
go test -race ./...

# Запуск с покрытием кода
go test -cover

# Генерация отчета о покрытии
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Запуск с таймаутом
go test -timeout 30s

# Параллельный запуск
go test -parallel 4

# Короткий режим (пропуск длинных тестов)
go test -short

# Запуск тестов N раз
go test -count=10

# Кэширование отключено
go test -count=1
```

### Флаги для контроля параллелизма

```go
func TestSomething(t *testing.T) {
    t.Parallel() // Тест будет запущен параллельно

    // test code
}
```

### Пропуск тестов

```go
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode")
    }

    // длительный тест
}

func TestIntegration(t *testing.T) {
    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration test")
    }

    // integration test
}
```

### Build tags

```go
// +build integration

package mypackage

import "testing"

func TestIntegration(t *testing.T) {
    // integration test
}
```

```bash
# Запуск с тегом
go test -tags=integration
```

## Бенчмарки

Бенчмарки используются для измерения производительности кода.

### Простой бенчмарк

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

### Что такое b.N?

`b.N` - это количество итераций, которое Go автоматически подбирает для получения стабильных результатов.

**Как это работает:**

1. Go начинает с `b.N = 1` и запускает бенчмарк
2. Если бенчмарк выполнился слишком быстро (< 1 секунды), Go увеличивает `b.N` (удваивает или увеличивает в 10 раз)
3. Процесс повторяется до тех пор, пока бенчмарк не будет выполняться достаточно долго для стабильных измерений
4. Go запускает бенчмарк несколько раз и вычисляет среднее время на операцию

```go
func BenchmarkExample(b *testing.B) {
    // Setup код выполняется 1 раз
    data := setupData()

    // Этот цикл выполнится b.N раз
    // b.N может быть: 1, 10, 100, 1000, 10000, 100000, 1000000...
    for i := 0; i < b.N; i++ {
        // Этот код будет измеряться
        processData(data)
    }
}
```

**Важные моменты:**

- **НЕ используйте b.N внутри логики**: `b.N` нужен только для цикла
- **Не изменяйте b.N**: Go сам подбирает оптимальное значение
- **Используйте b.ResetTimer()**: чтобы исключить setup из измерений

```go
// ❌ НЕПРАВИЛЬНО
func BenchmarkWrong(b *testing.B) {
    for i := 0; i < b.N; i++ {
        data := make([]int, b.N) // НЕ ДЕЛАЙТЕ ТАК!
        processData(data)
    }
}

// ✅ ПРАВИЛЬНО
func BenchmarkCorrect(b *testing.B) {
    data := make([]int, 1000) // Setup
    b.ResetTimer()           // Сброс таймера после setup

    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

**Пример работы b.N:**

```bash
# При запуске бенчмарка вы можете видеть что-то вроде:
# BenchmarkAdd-8   	  # первый запуск с b.N = 1
# BenchmarkAdd-8   	  # b.N = 100
# BenchmarkAdd-8   	  # b.N = 10000
# BenchmarkAdd-8   	1000000000   0.5234 ns/op  # финальный результат
```

### Бенчмарк с подготовкой

```go
func BenchmarkComplexOperation(b *testing.B) {
    data := generateTestData(1000)

    b.ResetTimer() // Сброс таймера после setup

    for i := 0; i < b.N; i++ {
        processData(data)
    }
}
```

### Параллельный бенчмарк

```go
func BenchmarkParallelOperation(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ExpensiveOperation()
        }
    })
}
```

### Табличные бенчмарки

```go
func BenchmarkFibonacci(b *testing.B) {
    benchmarks := []struct {
        name string
        n    int
    }{
        {"Fib10", 10},
        {"Fib20", 20},
        {"Fib30", 30},
    }

    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fibonacci(bm.n)
            }
        })
    }
}
```

### Бенчмарки с аллокацией памяти

```go
func BenchmarkStringConcat(b *testing.B) {
    b.ReportAllocs() // Показать информацию о выделении памяти

    for i := 0; i < b.N; i++ {
        result := ""
        for j := 0; j < 100; j++ {
            result += "a"
        }
    }
}

func BenchmarkStringBuilder(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 100; j++ {
            builder.WriteString("a")
        }
        _ = builder.String()
    }
}
```

### Запуск бенчмарков

```bash
# Запуск всех бенчмарков
go test -bench=.

# Запуск конкретного бенчмарка
go test -bench=BenchmarkAdd

# С информацией о памяти
go test -bench=. -benchmem

# Установка времени выполнения
go test -bench=. -benchtime=10s

# Количество итераций
go test -bench=. -benchtime=1000x

# Сравнение результатов
go test -bench=. > old.txt
# внести изменения
go test -bench=. > new.txt
benchcmp old.txt new.txt
```

### Интерпретация результатов

```
BenchmarkStringConcat-8         1000000    1234 ns/op    800 B/op    10 allocs/op
```

- `BenchmarkStringConcat-8` - название бенчмарка, `-8` означает GOMAXPROCS=8
- `1000000` - количество итераций (b.N)
- `1234 ns/op` - среднее время одной операции в наносекундах
- `800 B/op` - байт выделено на одну операцию (с флагом `-benchmem`)
- `10 allocs/op` - количество выделений памяти на операцию

### Сравнение производительности

```
BenchmarkStringConcat-8     1000000    1234 ns/op    800 B/op    10 allocs/op
BenchmarkStringBuilder-8   10000000     123 ns/op     64 B/op     1 allocs/op
```

StringBuilder в ~10 раз быстрее и использует в ~12 раз меньше памяти.

### Профилирование

```bash
# CPU профиль
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Память профиль
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof

# Trace
go test -bench=. -trace=trace.out
go tool trace trace.out
```

## Тестирование main

Функция `main()` требует особого подхода к тестированию.

### Вариант 1: Выделение логики

```go
// main.go
package main

import (
    "fmt"
    "os"
)

func main() {
    os.Exit(run())
}

func run() int {
    result, err := businessLogic()
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        return 1
    }

    fmt.Println(result)
    return 0
}

func businessLogic() (string, error) {
    return "success", nil
}
```

```go
// main_test.go
package main

import "testing"

func TestRun(t *testing.T) {
    exitCode := run()
    if exitCode != 0 {
        t.Errorf("expected exit code 0, got %d", exitCode)
    }
}

func TestBusinessLogic(t *testing.T) {
    result, err := businessLogic()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result != "success" {
        t.Errorf("got %q, want %q", result, "success")
    }
}
```

### Вариант 2: Тестирование с реальным main

```go
// main.go
package main

import (
    "flag"
    "fmt"
    "io"
    "os"
)

var (
    stdout io.Writer = os.Stdout
    stderr io.Writer = os.Stderr
    args   []string  = os.Args[1:]
)

func main() {
    os.Exit(mainWithExitCode())
}

func mainWithExitCode() int {
    fs := flag.NewFlagSet("myapp", flag.ContinueOnError)
    fs.SetOutput(stderr)

    verbose := fs.Bool("v", false, "verbose output")

    if err := fs.Parse(args); err != nil {
        return 2
    }

    if *verbose {
        fmt.Fprintln(stdout, "Running in verbose mode")
    }

    fmt.Fprintln(stdout, "Hello, World!")
    return 0
}
```

```go
// main_test.go
package main

import (
    "bytes"
    "testing"
)

func TestMain(t *testing.T) {
    tests := []struct {
        name           string
        args           []string
        expectedOutput string
        expectedExit   int
    }{
        {
            name:           "normal run",
            args:           []string{},
            expectedOutput: "Hello, World!\n",
            expectedExit:   0,
        },
        {
            name:           "verbose mode",
            args:           []string{"-v"},
            expectedOutput: "Running in verbose mode\nHello, World!\n",
            expectedExit:   0,
        },
        {
            name:           "invalid flag",
            args:           []string{"-invalid"},
            expectedOutput: "",
            expectedExit:   2,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var outBuf, errBuf bytes.Buffer

            stdout = &outBuf
            stderr = &errBuf
            args = tt.args

            defer func() {
                stdout = os.Stdout
                stderr = os.Stderr
                args = os.Args[1:]
            }()

            exitCode := mainWithExitCode()

            if exitCode != tt.expectedExit {
                t.Errorf("exit code = %d; want %d", exitCode, tt.expectedExit)
            }

            if tt.expectedOutput != "" {
                got := outBuf.String()
                if got != tt.expectedOutput {
                    t.Errorf("output = %q; want %q", got, tt.expectedOutput)
                }
            }
        })
    }
}
```

### Вариант 3: Интеграционные тесты через exec

```go
// main_integration_test.go
// +build integration

package main_test

import (
    "os/exec"
    "strings"
    "testing"
)

func TestMainIntegration(t *testing.T) {
    cmd := exec.Command("go", "run", "main.go")

    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("command failed: %v\noutput: %s", err, output)
    }

    got := strings.TrimSpace(string(output))
    want := "Hello, World!"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestMainWithFlags(t *testing.T) {
    cmd := exec.Command("go", "run", "main.go", "-v")

    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("command failed: %v", err)
    }

    if !strings.Contains(string(output), "verbose") {
        t.Error("expected verbose output")
    }
}
```

## Фаззинг

Фаззинг (fuzzing) - это техника автоматического тестирования, которая подаёт на вход программы случайные или мутированные данные для поиска багов, паник и уязвимостей.

### Зачем нужен фаззинг

1. **Находит неожиданные баги** - случаи, о которых вы не подумали
2. **Обнаруживает граничные случаи** - edge cases, которые сложно предугадать
3. **Выявляет проблемы безопасности** - переполнения буфера, инъекции и т.д.
4. **Находит панику** - неожиданные крэши приложения
5. **Дополняет unit-тесты** - находит то, что пропустили в тестах

### Базовый пример фаззинга (Go 1.18+)

```go
// strings.go
package mystrings

import "unicode/utf8"

func Reverse(s string) string {
    b := []byte(s)
    for i := 0; i < len(b)/2; i++ {
        b[i], b[len(b)-1-i] = b[len(b)-1-i], b[i]
    }
    return string(b)
}
```

```go
// strings_test.go
package mystrings

import (
    "testing"
    "unicode/utf8"
)

func FuzzReverse(f *testing.F) {
    // Seed corpus - начальные тестовые данные
    testcases := []string{"Hello, world", " ", "!12345", ""}
    for _, tc := range testcases {
        f.Add(tc)
    }

    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)

        if orig != doubleRev {
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }

        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
        }
    })
}
```

### Правильная реализация Reverse для UTF-8

```go
func Reverse(s string) string {
    r := []rune(s)
    for i := 0; i < len(r)/2; i++ {
        r[i], r[len(r)-1-i] = r[len(r)-1-i], r[i]
    }
    return string(r)
}
```

### Фаззинг с несколькими параметрами

```go
func FuzzAdd(f *testing.F) {
    f.Add(2, 3)
    f.Add(0, 0)
    f.Add(-1, 1)

    f.Fuzz(func(t *testing.T, a, b int) {
        result := Add(a, b)

        // Проверяем свойства операции
        if Add(a, b) != Add(b, a) {
            t.Errorf("addition is not commutative: Add(%d,%d)=%d, Add(%d,%d)=%d",
                a, b, Add(a, b), b, a, Add(b, a))
        }

        if a > 0 && b > 0 && result < a {
            t.Errorf("overflow detected: Add(%d,%d)=%d", a, b, result)
        }
    })
}
```

### Фаззинг парсера

```go
package parser

import (
    "encoding/json"
    "testing"
)

type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

func ParseConfig(data []byte) (*Config, error) {
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

func FuzzParseConfig(f *testing.F) {
    f.Add([]byte(`{"host":"localhost","port":8080}`))
    f.Add([]byte(`{}`))
    f.Add([]byte(`{"host":"example.com"}`))

    f.Fuzz(func(t *testing.T, data []byte) {
        cfg, err := ParseConfig(data)

        // Не должно быть паники
        if err != nil {
            return
        }

        // Если парсинг успешен, данные должны быть валидными
        if cfg.Port < 0 || cfg.Port > 65535 {
            t.Errorf("invalid port: %d", cfg.Port)
        }
    })
}
```

### Фаззинг с проверкой инвариантов

```go
func FuzzDivide(f *testing.F) {
    f.Add(10.0, 2.0)
    f.Add(0.0, 1.0)
    f.Add(-10.0, 2.0)

    f.Fuzz(func(t *testing.T, a, b float64) {
        result, err := Divide(a, b)

        if b == 0 {
            if err == nil {
                t.Error("expected error for division by zero")
            }
            return
        }

        if err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        // Проверяем инвариант: (a / b) * b ≈ a
        if b != 0 {
            check := result * b
            diff := a - check
            if diff < -0.0001 || diff > 0.0001 {
                t.Errorf("invariant failed: (%f / %f) * %f = %f, want ≈%f",
                    a, b, b, check, a)
            }
        }
    })
}
```

### Запуск фаззинга

```bash
# Запуск фаззинга (будет работать пока не найдет ошибку или пока не остановите)
go test -fuzz=FuzzReverse

# Запуск с ограничением по времени
go test -fuzz=FuzzReverse -fuzztime=30s

# Запуск определенное количество итераций
go test -fuzz=FuzzReverse -fuzztime=100000x

# Фаззинг с минимизацией найденных ошибок
go test -fuzz=FuzzReverse -fuzzminimizetime=1m

# Запуск обычных тестов + регрессионных фазз-тестов
go test

# Параллельные воркеры
go test -fuzz=FuzzReverse -parallel=8
```

### Результаты фаззинга

Когда фаззинг находит ошибку, он сохраняет входные данные:

```
testdata/fuzz/FuzzReverse/
    d5e8f2f3a1b2c4d6e7f8...
```

Эти файлы становятся частью регрессионных тестов и запускаются при обычном `go test`.

### Best practices для фаззинга

1. **Начните с хорошего seed corpus** - базовые тестовые случаи
2. **Проверяйте инварианты** - свойства, которые всегда должны выполняться
3. **Не игнорируйте ошибки** - даже ожидаемые ошибки должны обрабатываться корректно
4. **Используйте для критического кода** - парсеры, кодеки, криптография, безопасность
5. **Запускайте длительно** - дайте фаззеру время найти сложные баги
6. **Сохраняйте найденные кейсы** - используйте testdata/fuzz как регрессионные тесты

### Пример фаззинга URL парсера

```go
func FuzzParseURL(f *testing.F) {
    seeds := []string{
        "https://example.com",
        "http://user:pass@host:8080/path?key=value#fragment",
        "ftp://example.com",
        "//example.com",
    }

    for _, seed := range seeds {
        f.Add(seed)
    }

    f.Fuzz(func(t *testing.T, urlStr string) {
        u, err := url.Parse(urlStr)
        if err != nil {
            return
        }

        // После парсинга и сериализации должен получиться валидный URL
        serialized := u.String()
        _, err = url.Parse(serialized)
        if err != nil {
            t.Errorf("Parse->String->Parse failed: original=%q, serialized=%q, error=%v",
                urlStr, serialized, err)
        }
    })
}
```

## Продвинутые техники

### Использование testify

Библиотека testify предоставляет удобные функции для написания тестов.

**См. примеры в:** `examples/06_advanced/`

```bash
# Установка
go get github.com/stretchr/testify

# Запуск примеров
go test -v ./examples/06_advanced/
```

### Моки (Mocks)

Моки позволяют тестировать код, который зависит от внешних сервисов или компонентов.

**Пример из:** `examples/06_advanced/user_service_test.go`

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    expectedUser := &User{ID: 1, Name: "John", Email: "john@example.com"}

    mockRepo.On("GetByID", 1).Return(expectedUser, nil)

    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)

    require.NoError(t, err)
    assert.Equal(t, expectedUser, user)

    mockRepo.AssertExpectations(t)
}
```

### Табличные тесты для сложных сценариев

**Пример из:** `examples/06_advanced/user_service_test.go:32`

```go
tests := []struct {
    name        string
    user        *User
    setupMock   func(*MockUserRepository)
    wantErr     bool
    errContains string
}{
    {
        name: "success",
        user: &User{Name: "Alice", Email: "alice@example.com", Age: 25},
        setupMock: func(m *MockUserRepository) {
            m.On("Create", mock.AnythingOfType("*advanced.User")).Return(nil)
        },
        wantErr: false,
    },
    // ... другие тесты
}
```

## Дополнительные инструменты

### testify - популярная библиотека для тестирования

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWithTestify(t *testing.T) {
    result := Add(2, 3)

    // assert продолжает выполнение при ошибке
    assert.Equal(t, 5, result)
    assert.NotNil(t, result)

    // require останавливает выполнение при ошибке
    require.Equal(t, 5, result)
}
```

### Моки

```go
import (
    "testing"
    "github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetUser(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("GetUser", 1).Return(&User{ID: 1, Name: "John"}, nil)

    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)

    mockRepo.AssertExpectations(t)
}
```

### Golden files

```go
func TestRender(t *testing.T) {
    result := RenderHTML(&Data{Title: "Test", Content: "Hello"})

    golden := filepath.Join("testdata", "render.golden")

    if *update {
        os.WriteFile(golden, []byte(result), 0644)
    }

    expected, _ := os.ReadFile(golden)
    assert.Equal(t, string(expected), result)
}
```

## Директория testdata

`testdata` - это специальная директория в Go, которая игнорируется компилятором и предназначена для хранения тестовых данных.

### Особенности testdata

1. **Игнорируется go build**: Файлы в `testdata` не компилируются и не попадают в финальный бинарник
2. **Доступна в тестах**: Можно загружать файлы из `testdata` во время выполнения тестов
3. **Версионируется в Git**: Тестовые данные сохраняются в репозитории
4. **Стандартная практика**: Все Go разработчики знают это соглашение

### Структура testdata

```
mypackage/
├── mycode.go
├── mycode_test.go
└── testdata/
    ├── input.json           # Входные данные
    ├── output.json          # Ожидаемые результаты
    ├── golden/              # Golden files
    │   ├── result1.golden
    │   └── result2.golden
    └── fixtures/            # Фикстуры
        ├── user1.json
        └── user2.json
```

### Пример использования testdata

**См. примеры в:** `examples/07_testdata/`

```bash
# Запуск примеров
go test -v ./examples/07_testdata/
```

#### Чтение файлов из testdata

```go
func TestParseConfig(t *testing.T) {
    // Путь относительно текущего пакета
    data, err := os.ReadFile("testdata/config.json")
    if err != nil {
        t.Fatalf("failed to read test data: %v", err)
    }

    config, err := ParseConfig(data)
    if err != nil {
        t.Fatalf("ParseConfig failed: %v", err)
    }

    assert.Equal(t, "localhost", config.Host)
    assert.Equal(t, 8080, config.Port)
}
```

#### Табличные тесты с testdata

```go
func TestParseMultipleConfigs(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        wantHost string
        wantPort int
        wantErr  bool
    }{
        {
            name:     "valid config",
            filename: "testdata/valid.json",
            wantHost: "localhost",
            wantPort: 8080,
            wantErr:  false,
        },
        {
            name:     "invalid config",
            filename: "testdata/invalid.json",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            data, err := os.ReadFile(tt.filename)
            if err != nil {
                t.Fatalf("failed to read %s: %v", tt.filename, err)
            }

            config, err := ParseConfig(data)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.wantHost, config.Host)
            assert.Equal(t, tt.wantPort, config.Port)
        })
    }
}
```

#### Golden files в testdata

Golden files - это эталонные файлы, которые содержат ожидаемый результат работы функции.

```go
func TestRenderHTML(t *testing.T) {
    tests := []struct {
        name  string
        input *PageData
        golden string
    }{
        {
            name:   "simple page",
            input:  &PageData{Title: "Test", Body: "Hello World"},
            golden: "testdata/golden/simple.html",
        },
        {
            name:   "page with list",
            input:  &PageData{Title: "List", Items: []string{"A", "B", "C"}},
            golden: "testdata/golden/list.html",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := RenderHTML(tt.input)

            // Флаг для обновления golden files
            if *updateGolden {
                err := os.WriteFile(tt.golden, []byte(result), 0644)
                require.NoError(t, err)
                t.Log("Updated golden file:", tt.golden)
                return
            }

            // Сравнение с golden file
            expected, err := os.ReadFile(tt.golden)
            require.NoError(t, err)

            assert.Equal(t, string(expected), result)
        })
    }
}

// Флаг для обновления golden files
var updateGolden = flag.Bool("update", false, "update golden files")
```

**Обновление golden files:**

```bash
# Обновить все golden files
go test -update

# После обновления проверить что тесты проходят
go test
```

#### Использование embed для testdata

С Go 1.16+ можно встраивать testdata в бинарник с помощью `embed`:

```go
import _ "embed"

//go:embed testdata/config.json
var testConfigData []byte

func TestWithEmbed(t *testing.T) {
    config, err := ParseConfig(testConfigData)
    require.NoError(t, err)
    assert.Equal(t, "localhost", config.Host)
}

// Или для множества файлов
//go:embed testdata/*.json
var testFS embed.FS

func TestMultipleWithEmbed(t *testing.T) {
    data, err := testFS.ReadFile("testdata/config.json")
    require.NoError(t, err)
    // ...
}
```

#### Фикстуры в testdata

Фикстуры - это предопределенные данные для тестов.

```go
// testdata/fixtures/users.json
[
    {"id": 1, "name": "Alice", "email": "alice@example.com"},
    {"id": 2, "name": "Bob", "email": "bob@example.com"}
]

func loadUserFixtures(t *testing.T) []User {
    t.Helper()

    data, err := os.ReadFile("testdata/fixtures/users.json")
    require.NoError(t, err)

    var users []User
    err = json.Unmarshal(data, &users)
    require.NoError(t, err)

    return users
}

func TestUserService_WithFixtures(t *testing.T) {
    users := loadUserFixtures(t)

    for _, user := range users {
        t.Run(user.Name, func(t *testing.T) {
            valid := ValidateUser(&user)
            assert.True(t, valid)
        })
    }
}
```

### Best practices для testdata

1. **Организация**: Группируйте файлы по категориям (golden/, fixtures/, inputs/)
2. **Именование**: Используйте понятные имена файлов (valid_config.json, invalid_email.json)
3. **Размер**: Не храните слишком большие файлы (используйте минимальные данные)
4. **Формат**: Используйте текстовые форматы (JSON, XML, YAML) когда возможно
5. **Версионирование**: Коммитьте testdata в Git
6. **Документация**: Добавляйте README.md в testdata/ если структура сложная
7. **Очистка**: Регулярно проверяйте и удаляйте неиспользуемые файлы

### Пример структуры testdata

```
testdata/
├── README.md                  # Описание тестовых данных
├── fixtures/                  # Фикстуры для тестов
│   ├── users.json
│   ├── products.json
│   └── orders.json
├── golden/                    # Эталонные результаты
│   ├── report_monthly.html
│   ├── report_yearly.html
│   └── invoice.pdf
├── inputs/                    # Входные данные
│   ├── valid/
│   │   ├── config1.yaml
│   │   └── config2.yaml
│   └── invalid/
│       ├── missing_field.yaml
│       └── wrong_type.yaml
└── fuzz/                      # Данные для фаззинга (автогенерируется)
    └── FuzzParseConfig/
        └── ...
```

## Лучшие практики

1. **Именование тестов**
   - `TestFunctionName` для тестирования функции
   - `TestFunctionName_Scenario` для конкретного сценария
   - Используйте `t.Run()` для группировки

2. **Arrange-Act-Assert**
   - Подготовка данных
   - Выполнение действия
   - Проверка результата

3. **Изоляция**
   - Тесты не должны зависеть друг от друга
   - Используйте `t.Parallel()` для независимых тестов

4. **Читаемость**
   - Понятные имена переменных
   - Ясные сообщения об ошибках
   - Табличные тесты для множества случаев

5. **Coverage**
   - Стремитесь к 80%+ покрытию
   - Но не гонитесь за 100% - важно качество тестов

6. **Быстрые тесты**
   - Длительные тесты помечайте и пропускайте в `-short` режиме
   - Используйте моки для внешних зависимостей

7. **Тестируйте поведение, а не реализацию**
   - Фокусируйтесь на внешнем API
   - Не тестируйте приватные функции напрямую

## Полезные ссылки

- [Testing package](https://pkg.go.dev/testing)
- [Go Blog: Table Driven Tests](https://go.dev/blog/subtests)
- [Go Blog: Fuzzing](https://go.dev/doc/fuzz/)
- [Testify](https://github.com/stretchr/testify)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
