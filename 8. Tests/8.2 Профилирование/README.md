# Профилирование в Go

Примеры профилирования Go приложений для анализа производительности.

## 📁 Структура проекта

```
8.2 Профилирование/
├── 1. cpu profiling/         # CPU профилирование
│   ├── main.go              # Вычислительно-интенсивные операции
│   ├── go.mod
│   └── README.md
├── 2. memory profiling/      # Memory профилирование
│   ├── main.go              # Операции с памятью
│   ├── go.mod
│   └── README.md
├── 3. goroutine profiling/   # Goroutine профилирование
│   ├── main.go              # Конкурентные паттерны
│   ├── go.mod
│   └── README.md
├── 4. benchmark profiling/   # Benchmark профилирование
│   ├── operations.go        # Функции для бенчмарков
│   ├── benchmark_test.go    # Бенчмарки с профилированием
│   ├── go.mod
│   └── README.md
├── 5. http profiling/        # HTTP профилирование
│   ├── main.go              # HTTP сервер с pprof
│   ├── go.mod
│   └── README.md
└── README.md                # Общее руководство
```

## 🎯 Типы профилирования

### 1. 🔥 CPU Profiling
**Анализ использования процессора**

```bash
cd "1. cpu profiling"
go run main.go

# Анализ
go tool pprof cpu.prof
go tool pprof -http=:8080 cpu.prof
```

**Показывает:**
- Функции, потребляющие больше всего CPU времени
- Горячие пути в коде
- Bottlenecks в алгоритмах

### 2. 🧠 Memory Profiling  
**Анализ использования памяти**

```bash
cd "2. memory profiling"
go run main.go

# Анализ
go tool pprof mem.prof
go tool pprof -http=:8080 mem.prof
```

**Показывает:**
- Функции, аллоцирующие больше всего памяти
- Memory leaks
- Heap usage patterns

### 3. 🔄 Goroutine Profiling
**Анализ горутин**

```bash
cd "3. goroutine profiling"
go run main.go

# Анализ
go tool pprof goroutine.prof
go tool pprof -http=:8080 goroutine.prof
```

**Показывает:**
- Количество активных горутин
- Блокировки и ожидания
- Goroutine leaks

### 4. ⚡ Benchmark Profiling
**Производительность через бенчмарки**

```bash
cd "4. benchmark profiling"

# Обычные бенчмарки
go test -bench=.

# С CPU профилем
go test -bench=. -cpuprofile=cpu.prof

# С Memory профилем
go test -bench=. -memprofile=mem.prof

# С Mutex профилем
go test -bench=. -mutexprofile=mutex.prof
```

### 5. 🌐 HTTP Profiling
**Профилирование HTTP сервера**

```bash
cd "5. http profiling"
go run main.go

# В браузере:
# http://localhost:8080/debug/pprof/
# http://localhost:8080/debug/pprof/profile
# http://localhost:8080/debug/pprof/heap
```

## 🛠 Основные команды

### Сбор профилей
```bash
# CPU профиль
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Memory профиль
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine профиль
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

### Анализ в консоли
```bash
go tool pprof profile.prof
(pprof) top
(pprof) list main.slowFunction
(pprof) web
```

### Веб-интерфейс
```bash
go tool pprof -http=:8080 profile.prof
```

## 📊 Что анализировать

### CPU Profiling
- `top` - топ функций по CPU времени
- `list funcName` - код функции с временами
- `web` - граф вызовов

### Memory Profiling
- `top` - топ функций по аллокациям
- `list funcName` - где происходят аллокации
- `inuse_space` vs `alloc_space`

### Goroutine Profiling
- Количество горутин
- Состояния горутин (running, waiting, blocked)
- Stack traces заблокированных горутин

## 🎯 Советы по оптимизации

### CPU
- Избегайте ненужных вычислений в циклах
- Используйте более эффективные алгоритмы
- Кэшируйте результаты дорогих операций

### Memory
- Переиспользуйте объекты (sync.Pool)
- Избегайте частых аллокаций в горячих путях
- Используйте strings.Builder вместо конкатенации строк

### Goroutines
- Не создавайте горутины без контроля
- Используйте context для отмены
- Правильно закрывайте каналы

## 🚀 Быстрый старт

1. Выберите тип профилирования
2. Перейдите в соответствующую папку
3. Запустите пример
4. Проанализируйте результаты
5. Экспериментируйте с кодом