# Memory Profiling

Пример Memory профилирования для анализа использования памяти.

## Запуск

```bash
go run main.go
```

Программа создаст файл `mem.prof` с данными профилирования памяти.

## Анализ в консоли

```bash
# Анализ текущих аллокаций
go tool pprof mem.prof

# Анализ всех аллокаций за время работы
go tool pprof -alloc_space mem.prof
```

### Основные команды в pprof:
```
(pprof) top                    # топ функций по памяти
(pprof) top -cum              # с учетом вызываемых функций
(pprof) list main.createManySlices
(pprof) web                   # граф использования памяти
(pprof) png                   # сохранить как PNG
```

### Фильтрация данных:
```
(pprof) top --cum main        # только функции из main пакета
(pprof) list -cum createLarge # детали конкретной функции
```

## Анализ в браузере

```bash
# Веб-интерфейс для анализа памяти
go tool pprof -http=:8080 mem.prof

# Анализ всех аллокаций
go tool pprof -http=:8080 -alloc_space mem.prof
```

### Веб-интерфейс показывает:
- **Graph** - граф аллокаций памяти
- **Flame Graph** - flame график использования памяти
- **Top** - функции, потребляющие больше всего памяти
- **Source** - аннотированный код с метриками

## Что анализирует профиль

1. **createManySlices** - множественные slice аллокации
2. **createLargeStructures** - большие структуры с вложенными данными
3. **simulateMemoryLeaks** - накопление данных в глобальных переменных

## Типы профилей памяти

### По умолчанию (inuse):
- **inuse_space** - сколько памяти используется сейчас
- **inuse_objects** - количество объектов в памяти

### Все аллокации (-alloc_space):
- **alloc_space** - общий объем аллокаций
- **alloc_objects** - общее количество созданных объектов

## Полезные команды

```bash
# Сравнение двух профилей
go tool pprof -base=old.prof new.prof

# Профилирование работающего приложения
go tool pprof http://localhost:6060/debug/pprof/heap

# Экспорт в различные форматы
go tool pprof -top mem.prof
go tool pprof -list=main.createManySlices mem.prof
go tool pprof -web mem.prof
```

## Интерпретация результатов

- **flat** - память, аллоцированная непосредственно функцией
- **cum** - общая память (включая вызываемые функции)
- Высокое значение **cum** указывает на функции, вызывающие много аллокаций
- Высокое значение **flat** указывает на функции, напрямую создающие объекты