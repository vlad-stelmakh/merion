# Benchmark Profiling

Примеры профилирования через бенчмарки Go для измерения производительности.

## 🚀 Быстрый старт - 5 основных команд

### 1. Запуск базового бенчмарка с метриками памяти
```bash
go test -bench=. -benchmem
```

### 2. Создание CPU профиля и запуск веб-анализа
```bash
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -http=:8080 cpu.prof
```

### 3. Создание Memory профиля и запуск веб-анализа
```bash
go test -bench=. -memprofile=mem.prof && go tool pprof -http=:8080 mem.prof
```

### 4. Полное профилирование с сохранением результатов
```bash
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem > benchmark_results.txt
```

### 5. Интерактивный анализ CPU профиля в консоли
```bash
go test -bench=. -cpuprofile=cpu.prof && go tool pprof cpu.prof
```

---

## Запуск бенчмарков

### Базовый запуск:
```bash
go test -bench=.
go test -bench=. -benchmem
```

### Детальный анализ:
```bash
# Запуск с профилированием CPU
go test -bench=. -cpuprofile=cpu.prof

# Запуск с профилированием памяти
go test -bench=. -memprofile=mem.prof

# Комбинированное профилирование
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem
```

## Анализ профилей в консоли

### CPU профиль:
```bash
go tool pprof cpu.prof
(pprof) top
(pprof) list BenchmarkStringConcat
(pprof) web
```

### Memory профиль:
```bash
go tool pprof mem.prof
(pprof) top
(pprof) list -cum AllocateMany
(pprof) traces
```

## Анализ в браузере

```bash
# CPU анализ
go tool pprof -http=:8080 cpu.prof

# Memory анализ  
go tool pprof -http=:8080 mem.prof

# Объединенный анализ (если есть оба профиля)
go tool pprof -http=:8080 -base=cpu.prof mem.prof
```

## Конкретные бенчмарки

### 1. Сравнение конкатенации строк:
```bash
go test -bench=BenchmarkStringConcat -benchmem
```

### 2. Сравнение алгоритмов сортировки:
```bash
go test -bench=BenchmarkSorting -benchmem
```

### 3. Анализ аллокаций памяти:
```bash
go test -bench=BenchmarkAllocation -benchmem
```

### 4. Value vs Pointer методы:
```bash
go test -bench=BenchmarkPointDistance -benchmem
```

## Дополнительные опции

### Контроль времени выполнения:
```bash
# Минимум 10 секунд на бенчмарк
go test -bench=. -benchtime=10s

# Конкретное количество итераций
go test -bench=. -benchtime=1000000x
```

### Фильтрация бенчмарков:
```bash
# Только бенчмарки строк
go test -bench=.*String.*

# Исключить медленные тесты
go test -bench=. -short
```

### Сравнение результатов:
```bash
# Сохранить результат в файл
go test -bench=. > old.txt

# После изменений
go test -bench=. > new.txt

# Сравнить (требует benchcmp)
benchcmp old.txt new.txt
```

## Интерпретация результатов

### Пример вывода:
```
BenchmarkStringConcat/Slow-8    100000   12345 ns/op   1024 B/op   10 allocs/op
BenchmarkStringConcat/Fast-8    500000    2468 ns/op    256 B/op    1 allocs/op
```

### Объяснение полей:
- **100000** - количество итераций
- **12345 ns/op** - наносекунд на операцию
- **1024 B/op** - байт аллокаций на операцию
- **10 allocs/op** - количество аллокаций на операцию
- **-8** - количество GOMAXPROCS

## Профилирование конкретных функций

### CPU hotspots:
```bash
go test -bench=BenchmarkSorting -cpuprofile=sort_cpu.prof
go tool pprof -http=:8080 sort_cpu.prof
# Анализируйте flame graph для поиска узких мест
```

### Memory usage:
```bash
go test -bench=BenchmarkAllocation -memprofile=alloc_mem.prof
go tool pprof -http=:8080 alloc_mem.prof
# Смотрите на allocation patterns
```

## Лучшие практики

### ✅ Do:
- Используйте `b.ResetTimer()` после инициализации
- Тестируйте на реальных данных
- Проверяйте и CPU, и память
- Сравнивайте альтернативные реализации

### ❌ Don't:
- Не включайте инициализацию в измерения
- Не игнорируйте garbage collector эффекты
- Не делайте преждевременную оптимизацию
- Не тестируйте только на маленьких данных

## Автоматизация

### Скрипт для полного анализа:
```bash
#!/bin/bash
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem > bench.txt
echo "Results saved to bench.txt"
echo "CPU profile: go tool pprof -http=:8080 cpu.prof"
echo "Memory profile: go tool pprof -http=:8080 mem.prof"
```

## CI/CD интеграция

```bash
# Проверка деградации производительности
go test -bench=. -count=5 | tee bench_results.txt
# Анализируйте результаты на предмет регрессии
```