# Отладка в Go

Практические примеры отладки Go кода для эффективного поиска и устранения багов.

## 📁 Структура проекта

```
8.3 Отладка/
├── 1. Bugs/             # Типичные баги и их поиск
│   ├── bugs.go         # Примеры различных багов
│   └── go.mod
├── 2. Concurrency/      # Отладка горутин и каналов
│   ├── concurrent.go   # Конкурентные паттерны
│   └── go.mod
└── README.md           # Инструкции по отладке
```

## 🔍 Основные техники отладки

### 1. **Breakpoints (Точки останова)**

#### Установка breakpoint:
- **Клик на полях слева** от номера строки
- Или **Ctrl+F8** (Windows/Linux) / **Cmd+F8** (macOS)

#### Типы breakpoints:
- **Обычный** - останавливается всегда
- **Условный** - останавливается при выполнении условия
- **Временный** - срабатывает один раз

### 2. **Пошаговое выполнение**

- **Step Over (F8)** - переход к следующей строке
- **Step Into (F7)** - заход внутрь функции
- **Step Out (Shift+F8)** - выход из текущей функции
- **Run to Cursor (Alt+F9)** - выполнение до курсора

### 3. **Исследование переменных**

- **Variables panel** - просмотр локальных переменных
- **Watches** - отслеживание конкретных выражений
- **Evaluate Expression (Alt+F8)** - выполнение кода в контексте

## 🎯 Примеры отладки

### 1. Bugs - Типичные ошибки

```bash
cd "1. Bugs"
go run bugs.go
```

**Демонстрирует:**
- Index out of bounds
- Логические ошибки
- Race conditions

### 2. Concurrency - Конкурентность

```bash
cd "2. Concurrency" 
go run concurrent.go
```

**Демонстрирует:**
- Отладку горутин
- Работу с каналами
- Race conditions и их предотвращение

## 🛠 Настройка отладки в IDE

### GoLand:
1. **Run → Edit Configurations...**
2. **+** → **Go Build**
3. Настройте файл и рабочую директорию
4. Используйте **Debug** режим

### VS Code:
1. Установите Go extension
2. Создайте `.vscode/launch.json`
3. Настройте конфигурацию для отладки

## 🚨 Отладка Race Conditions

### Детекция:
```bash
go run -race program.go
go test -race
```

### Основные признаки:
- Непредсказуемые результаты
- Sporadic failures
- Зависимость от timing

### Решения:
- Используйте **mutex** для синхронизации
- Применяйте **channels** для коммуникации
- Используйте **atomic** операции

## 🔧 Полезные команды

### Сборка с отладочной информацией:
```bash
go build -gcflags="-N -l" program.go
```

### Запуск с race detector:
```bash
go run -race program.go
```

### Profiling для отладки производительности:
```bash
go run program.go -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## 📝 Советы по отладке

### 1. **Systematic Approach**
- Воспроизведите баг
- Сократите код до минимума
- Используйте binary search для локализации

### 2. **Logging vs Debugging**
- Логи для production
- Debugger для development
- Структурированные логи (JSON)

### 3. **Common Patterns**
- Проверяйте nil pointers
- Валидируйте входные данные
- Обрабатывайте errors

### 4. **Concurrency Issues**
- Всегда используйте `-race` flag
- Минимизируйте shared state
- Предпочитайте channels над mutex

## 🎯 Практические упражнения

1. **Запустите примеры с breakpoints**
2. **Исследуйте переменные в разных состояниях**
3. **Попробуйте исправить баги**
4. **Добавьте свои примеры отладки**

## 📚 Дополнительные ресурсы

- [Delve Debugger](https://github.com/go-delve/delve)
- [Go Race Detector](https://golang.org/doc/articles/race_detector.html)
- [Debugging with GoLand](https://www.jetbrains.com/help/go/debugging-code.html)