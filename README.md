# Pantan 🚀

![CI/CD Pipeline](https://github.com/smukm/pantan/actions/workflows/main.yml/badge.svg)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green.svg)

**Pantan** — сервис с полностью автоматизированным CI/CD пайплайном, интегрированным стеком наблюдаемости (Observability) и умной генерацией changelog'ов с помощью искусственного интеллекта.

## 🌟 Особенности

- **Полный цикл CI/CD**: Автоматическая проверка кода, сборка multi-arch Docker образов (`amd64`/`arm64`) и zero-downtime деплой на сервер по SSH.
- **AI-генерация релизов**: Автоматическое создание профессиональных Release Notes на основе истории коммитов с использованием LLM (OpenRouter / Nvidia Nemotron).
- **Уведомления**: Мгновенные оповещения о статусе деплоя и релиза в Telegram.
- **Production-Ready Observability**: Встроенный полнофункциональный стек мониторинга и логирования (Prometheus, Grafana, Loki, Promtail, Node Exporter).
- **Безопасность**: Базовая аутентификация для Prometheus, автоматические проверки уязвимостей (`govulncheck`) и статический анализ (`golangci-lint`).
- **Умное логирование**: Глобальная конфигурация Docker-логирования с ротацией и сжатием для экономии дискового пространства.

## 🏗 Архитектура и стек

| Компонент | Технология | Описание |
|-----------|------------|----------|
| **Backend** | Go 1.25 | Основная бизнес-логика сервиса |
| **CI/CD** | GitHub Actions | Пайплайн из 5 шагов (Tests ➔ Build ➔ Deploy ➔ Release ➔ Notify) |
| **Контейнеризация** | Docker & Docker Compose | Изоляция и оркестрация сервисов |
| **Метрики** | Prometheus + Node Exporter | Сбор метрик приложения и хоста |
| **Логи** | Loki + Promtail | Агрегация и централизованный сбор логов (в т.ч. внешних сайтов) |
| **Визуализация** | Grafana | Дашборды для метрик и логов с автоматическим provisioning'ом |
| **AI** | OpenRouter API | Генерация описаний к релизам |

## ⚙️ Переменные окружения (`.env`)

Для работы проекта в production-окружении необходимо создать файл `.env` в директории деплоя на сервере (рядом с `docker-compose.prod.yml`).

```env
# === Приложение ===
PT_PORT=8033
PT_ENV=production
PT_LOG_LEVEL=1
PT_LOG_FORMAT=json
PT_LOG_TARGET=local

# === Grafana ===
GRAFANA_USER=admin
GRAFANA_PASSWORD=your_secure_password

# === Registry ===
REGISTRY=ghcr.io
IMAGE_NAME=smukm/pantan
```

### Локальный запуск
```bash
docker compose -f docker-compose.prod.yml up -d --build
```

### CI/CD Пайплайн
Этапы выполнения:
- Code Quality: Проверка форматирования (gofmt), статический анализ (golangci-lint), тесты с race-детектором и поиск уязвимостей (govulncheck).
- Build & Push: Сборка Docker образа с кэшированием и публикацией в GitHub Container Registry (GHCR).
- Deploy: Подключение к серверу по SSH, авторизация в GHCR, скачивание нового образа и перезапуск сервисов через docker compose.
- Create Release: Анализ истории коммитов и измененных файлов, генерация Release Notes через OpenRouter API и создание тега/релиза в GitHub.
- Notify: Отправка сообщения в Telegram с итогами деплоя и сгенерированным changelog'ом.

#### Секреты для GitHub Actions
Для работы CI/CD необходимо добавить следующие секреты в настройки репозитория (Settings ➔ Secrets and variables ➔ Actions):
- SSH_HOST, SSH_USERNAME, SSH_PASSWORD, SSH_PORT — для подключения к серверу.
- TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID — для уведомлений.
- OPENROUTER_API_KEY — для генерации Release Notes через ИИ.
#### Observability (Мониторинг и Логи)
Все сервисы мониторинга поднимаются вместе с приложением и изолированы в единой сети.
- Grafana: http://<your-server-ip>:3000
Логин: admin | Пароль: задается в .env
Преднастроенные datasource'ы для Prometheus и Loki.
- Prometheus: http://<your-server-ip>:9090
Защищен базовой аутентификацией (файл web.yml).
Хранение метрик: 7 дней (--storage.tsdb.retention.time=7d).
- Loki: Внутренний порт 3100 (используется Grafana для запросов логов).
- Promtail: Автоматически парсит логи сайта из /var/www/www-root/data/www/brend-logo.ru и отправляет их в Loki.
- Логирование Docker: Глобально настроен json-file драйвер с ротацией (макс. 5 файлов по 5MB, сжатие включено), что предотвращает переполнение диска логами контейнеров.

### Структура конфигураций на сервере
/var/www/www-root/data/stats/
├── docker-compose.prod.yml   # Оркестрация сервисов
├── .env                      # Переменные окружения
└── Docker/                   # Конфигурации для observability стека
├── prometheus/
│   ├── prometheus-prod.yml
│   └── web.yml           # Конфиг базовой аутентификации
├── grafana/
│   └── provisioning/     # Автоматическое добавление Data Sources
├── loki/
│   └── loki-config.yml
└── promtail/
└── promtail-config.yml