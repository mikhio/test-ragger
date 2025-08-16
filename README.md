# test-ragger

🔍 **RAG система для поиска по HTML документам** с использованием векторных эмбеддингов OpenAI и Qdrant.

## 🚀 Быстрый старт

```bash
# 1. Настройка проекта
make setup

# 2. Запуск Qdrant базы данных
make docker-up

# 3. Установка OpenAI API ключа
export OPENAI_API_KEY=sk-your-key-here

# 4. Индексация HTML файлов
make ingest

# 5. Поиск по документам
make search Q="ваш поисковый запрос"
```

## 📖 Документация

- **[📋 Makefile команды](docs/MAKEFILE_CHEATSHEET.md)** - Полный список команд и примеры
- **[🛠️ Локальная разработка](docs/LOCAL_DEVELOPMENT.md)** - Подробная настройка среды разработки
- **[🔧 Chunker утилита](docs/chunker.md)** - Документация по разбиению текста на чанки

## 🏗️ Архитектура

```
cmd/test-ragger/           # Точка входа приложения
internal/
├── configure/            # DI контейнер и конфигурация
├── usecase/             # Бизнес-логика (ingest, search)
├── models/              # Модели данных
└── utils/               # Утилиты (chunker, htmlx, prompt)
```

## ⚡ Основные команды

```bash
make help           # Показать все доступные команды
make ingest         # Индексация HTML файлов
make search Q="..."  # Поиск по индексированным данным
make docker-up      # Запуск Qdrant
make docker-down    # Остановка Qdrant
```

## 🔧 Технологии

- **Go 1.21+** - основной язык
- **OpenAI API** - создание эмбеддингов
- **Qdrant** - векторная база данных
- **Docker** - контейнеризация Qdrant
- **Make** - автоматизация команд

## 📊 Веб-интерфейс Qdrant

После запуска `make docker-up` доступен по адресу:
- 🌐 **Dashboard**: http://localhost:6333/dashboard
- 📚 **API docs**: http://localhost:6333/docs

## ❓ Troubleshooting

Если возникают проблемы:
1. Проверьте [документацию по разработке](docs/LOCAL_DEVELOPMENT.md)
2. Используйте `make help` для просмотра команд
3. Проверьте логи Qdrant: `make docker-logs`

---
*Создано для эффективного семантического поиска по документации*