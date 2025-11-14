# Команды для профилирования

## Быстрый запуск примеров

```bash
# 1. CPU profiling
cd "1. cpu profiling"
go run main.go
go tool pprof -http=:8080 cpu.prof

# 2. Memory profiling  
cd "2. memory profiling"
go run main.go
go tool pprof -http=:8080 mem.prof

# 3. Goroutine profiling
cd "3. goroutine profiling"
go run main.go
go tool pprof -http=:8080 goroutine.prof

# 4. Benchmark profiling
cd "4. benchmark profiling"
go test -bench=. -cpuprofile=cpu.prof
go tool pprof -http=:8080 cpu.prof

# 5. HTTP profiling
cd "5. http profiling"
go run main.go
# Открыть: http://localhost:8080/debug/pprof/
```

## Анализ профилей

### В консоли
```bash
go tool pprof profile.prof
(pprof) top           # Топ функций
(pprof) list main     # Код main функции
(pprof) web           # Граф в браузере
(pprof) exit
```

### В браузере
```bash
go tool pprof -http=:8080 profile.prof
```

## Бенчмарки с профилированием

```bash
# CPU профиль
go test -bench=. -cpuprofile=cpu.prof

# Memory профиль
go test -bench=. -memprofile=mem.prof

# Mutex профиль
go test -bench=. -mutexprofile=mutex.prof

# Block профиль
go test -bench=. -blockprofile=block.prof
```

## HTTP профилирование

```bash
# Сбор профилей с живого сервера
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/goroutine
```