# Команды для отладки

## Быстрый запуск примеров

```bash
# 1. Bugs - Типичные ошибки
cd "1. Bugs"
go run bugs.go

# 2. Concurrency - Конкурентные паттерны  
cd "2. Concurrency"
go run concurrent.go
```

## Отладка с race detector

```bash
# Обнаружение race conditions
go run -race bugs.go
go run -race concurrent.go
```

## Сборка с отладочной информацией

```bash
# Полная отладочная информация
go build -gcflags="-N -l" bugs.go
go build -gcflags="-N -l" concurrent.go
```

## Debugging в разных IDE

### GoLand
1. **Run → Edit Configurations**
2. Выберите файл для отладки
3. Установите breakpoints (клик слева от строки)
4. **Debug** вместо **Run**

### VS Code
1. **F5** для запуска отладки
2. Выберите конфигурацию Go
3. Установите breakpoints
4. Используйте панель Variables

## Полезные команды отладки

```bash
# Trace execution
go run -trace=trace.out program.go
go tool trace trace.out

# Memory profiling
go run -memprofile=mem.prof program.go  
go tool pprof mem.prof

# CPU profiling
go run -cpuprofile=cpu.prof program.go
go tool pprof cpu.prof
```

## Основные hotkeys

- **F8** - Step Over (следующая строка)
- **F7** - Step Into (заход в функцию)  
- **Shift+F8** - Step Out (выход из функции)
- **F9** - Resume (продолжить выполнение)
- **Ctrl+F8** - Toggle Breakpoint