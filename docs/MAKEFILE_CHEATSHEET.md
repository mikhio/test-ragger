# 📋 Makefile Шпаргалка

> [← Назад к документации](README.md) | [🏠 Главная](../README.md)

## 🚀 Быстрый старт

```bash
make setup                          # Первоначальная настройка
make docker-up                      # Запуск Qdrant
export OPENAI_API_KEY=sk-...        # Установка API ключа
make ingest                         # Индексация HTML файлов
make search Q="машинное обучение"   # Поиск
```

## 📖 Основные команды

### Настройка и сборка
```bash
make setup        # Настройка проекта + загрузка зависимостей
make build        # Сборка приложения
make build-clean  # Чистая сборка (с очисткой)
make deps         # Только загрузка зависимостей
```

### Запуск приложения
```bash
# Индексация
make ingest                              # Базовая индексация
make ingest DIR=./my-docs               # Кастомная папка
make ingest MODEL=text-embedding-3-large # Другая модель

# Поиск
make search Q="запрос"                   # Базовый поиск
make search Q="AI" K=10                  # Топ-10 результатов
make search Q="ML" K=5 LANG=ru          # С фильтром языка
```

### Docker управление
```bash
make docker-up      # Запуск Qdrant
make docker-down    # Остановка Qdrant
make docker-logs    # Просмотр логов в реальном времени
make docker-clean   # Остановка + очистка данных
```

### Разработка
```bash
make test           # Запуск тестов
make lint           # Форматирование + проверка кода
make fmt            # Только форматирование
make vet            # Только проверка кода
make tidy           # Очистка go.mod
```

### Очистка
```bash
make clean          # Удаление бинарников
make clean-all      # Полная очистка (+ Docker данные + mod cache)
```

## 🔧 Переменные окружения

```bash
# Обязательная
export OPENAI_API_KEY=sk-your-key

# Опциональные (с значениями по умолчанию)
DIR=./html                    # Папка с HTML файлами
MODEL=text-embedding-3-small  # Модель эмбеддингов
QDRANT=localhost:6334         # Адрес Qdrant gRPC
K=5                          # Количество результатов поиска
LANG=                        # Фильтр языка (пустой = все языки)
```

## 💡 Примеры использования

```bash
# Полный цикл разработки
make setup
make docker-up
export OPENAI_API_KEY=sk-...
make ingest DIR=./docs MODEL=text-embedding-3-large
make search Q="векторные базы данных" K=10

# Быстрая разработка
make build && ./bin/test-ragger -mode=search -q="test" -k=3

# Отладка проблем
make docker-logs
make lint
make clean-all && make setup
```

## 🎯 Полезные фишки

1. **Автоматический GOPROXY**: Makefile автоматически использует публичный прокси
2. **Умная сборка**: Создает папки автоматически
3. **Проверка параметров**: Проверяет обязательные параметры (например, Q для search)
4. **Эмодзи и цвета**: Удобный вывод с иконками
5. **Справка**: `make help` покажет все доступные команды
