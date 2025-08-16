# 🛠️ Локальная разработка с Docker Qdrant

> [← Назад к документации](README.md) | [🏠 Главная](../README.md)

## Настройка окружения

### 1. Запуск Docker Desktop
Убедитесь, что Docker Desktop запущен на вашем Mac.

### 2. Создание .env файла
```bash
cp env-example .env
```

Отредактируйте `.env` файл и добавьте ваш OpenAI API ключ:
```bash
OPENAI_API_KEY=sk-your-actual-api-key-here
```

### 3. Запуск Qdrant базы данных
```bash
# Запуск Qdrant в фоне
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs qdrant
```

### 4. Проверка подключения к Qdrant
- **Веб-интерфейс**: http://localhost:6333/dashboard
- **API документация**: http://localhost:6333/docs
- **gRPC endpoint**: localhost:6334 (для приложения)

## Разработка и отладка

### Сборка приложения
```bash
go build -o bin/test-ragger ./cmd/test-ragger
```

### Индексация тестовых данных
```bash
# Использование примера HTML файла
./bin/test-ragger -mode=ingest -dir=./html

# Или с go run для отладки
go run ./cmd/test-ragger -mode=ingest -dir=./html
```

### Поиск по индексированным данным
```bash
# Поиск с базовыми параметрами
./bin/test-ragger -mode=search -q="машинное обучение" -k=5

# Поиск с дополнительными параметрами
./bin/test-ragger -mode=search -q="нейронные сети" -k=10 -lang=ru

# С go run для отладки
go run ./cmd/test-ragger -mode=search -q="векторные базы данных" -k=3
```

## Отладка в IDE

### VS Code
1. Создайте `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Ingest",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/test-ragger",
            "args": ["-mode=ingest", "-dir=./html"],
            "env": {}
        },
        {
            "name": "Debug Search",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/test-ragger",
            "args": ["-mode=search", "-q=машинное обучение", "-k=5"],
            "env": {}
        }
    ]
}
```

2. Поставьте брейкпоинты в коде
3. Нажмите F5 для запуска отладки

### GoLand/IntelliJ
1. Создайте Run Configuration
2. Program: `cmd/test-ragger`
3. Arguments: `-mode=search -q="ваш запрос" -k=5`
4. Working directory: корень проекта

## Полезные команды

### Docker управление
```bash
# Остановка Qdrant
docker-compose down

# Перезапуск Qdrant
docker-compose restart

# Очистка данных (удаление всех векторов)
docker-compose down -v

# Просмотр логов в реальном времени
docker-compose logs -f qdrant
```

### Проверка состояния Qdrant
```bash
# Проверка через curl
curl http://localhost:6333/collections

# Информация о коллекции docs
curl http://localhost:6333/collections/docs
```

## Структура данных

После индексации в Qdrant создается коллекция `docs` с векторами размерностью 1536 (OpenAI embeddings).

Каждый вектор содержит payload:
- `title` - заголовок документа
- `text` - текст чанка
- `path` - путь к файлу
- `lang` - язык (если указан)

## Troubleshooting

### Проблемы с подключением к Qdrant
1. Проверьте, что Docker запущен: `docker ps`
2. Проверьте порты: `lsof -i :6333,6334`
3. Перезапустите Qdrant: `docker-compose restart qdrant`

### Проблемы с OpenAI API
1. Проверьте .env файл и API ключ
2. Проверьте баланс на счету OpenAI
3. Убедитесь, что модель `text-embedding-3-small` доступна

### Пустые результаты поиска
1. Убедитесь, что индексация прошла успешно
2. Проверьте коллекцию в веб-интерфейсе Qdrant
3. Попробуйте другие поисковые запросы
