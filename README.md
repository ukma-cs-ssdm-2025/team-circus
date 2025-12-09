# Markdown Circus Docs

[![CI](https://github.com/ukma-cs-ssdm-2025/team-circus/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ukma-cs-ssdm-2025/team-circus/actions/workflows/ci.yml)

Наш проект MCD (Markdown Circus Docs) - це інструмент спільного редагування документів у форматі Markdown.

**Чому ми обрали саме Markdown?**

- Легко та швидко редагується
- Має уніфікований і стандартизований синтаксис
- Підтримується більшістю IDE та редакторів
- Знайомий і зрозумілий кожному члену команди розробки

## Учасники команди

- Артур Ключка
  - github: [velosypedno](https://github.com/velosypedno)
- Владислав Панько
  - github: [xncff](https://github.com/xncff)
- Олексій Костик
  - github: [kitmelancholic](https://github.com/kitmelancholic)
- Ярослав Шах
  - github: [YaroslavShakh](https://github.com/YaroslavShakh)

Командний статут знаходиться [тут](./TeamChapter.md).

## Запуск проєкту

### Передумови

Потрібно мати встановлені:

- **Go** (версія 1.23+)
- **Bun** (версія 1.3+)
- **Task** (taskfile.dev)
- **Docker** (з docker-compose)

Клонування репозиторію

```bash
git clone https://github.com/ukma-cs-ssdm-2025/team-circus
cd team-circus
```

Скопіювати .env

```bash
task copy:env
```

### Postgres

1. **Підняти контейнер з бд**

   ```bash
   task docker:postgres:up
   ```

2. **Заранити міграції бд**

   ```bash
   task docker:migrator:up
   ```

### Backend

Доступні дві опції:

1. Запустити у контейнері
2. Запустити локально

#### У контейнері

1. **Підняти контейнер з бекендом**

   ```bash
   task docker:backend:up
   ```

#### Локально

1. **Встановити залежності**

   ```bash
   task back:download
   ```

2. **Додати залежнсоті у локальне оточення**

   ```bash
   task back:vendor
   ```

3. **Скопіювати модифікований .env**

   ```bash
   task back:copy:env
   ```

4. **Запустит бекенд локально**

   ```bash
   task back:run
   ```

### Frontend

1. **Встановлення залежностей**

   ```bash
   task front:install
   ```

2. **Запуск frontend**

   ```bash
   task front:dev
   ```

### Доступні команди

Для перегляду всіх доступних команд використовуйте:

```bash
task --list-all
```

## Тестування

### Юніт тести

```bash
task back:test:unit
```

Команда запусутить всі юніт тести, тобто тестові файли, які не мають жодних тегів по типу `func_test`. Також Буде виведеий coverage по кожному пакету.

#### Відображення coverage у інтерактивному форматі

```bash
task back:test:unit:coverage
```

Згенерує звіт, який буде лежати в [`./backend/coverage/coverage.out`](./backend/coverage/coverage.out)

```bash
task back:test:unit:coverage:html
```

Згенерує html, який буде лежати в [`./backend/coverage/coverage.html`](./backend/coverage/coverage.html)

### Функціональні (інтеграційні) тести

З функціональними тестами ситуація трохи складніша.
Ці тести часто потреюуть зовнішніх залежностей, у нашому випадку - бд.
Залежності можна додати в [`backend/docker-compose.test.yml`](./backend/docker-compose.test.yml).
Там же ми хардкодимо значення змінних оточення.
У [`backend/tests/pkg/testapp/app.go`](./backend/tests/pkg/testapp/app.go) в конфігу ми вказуємо всі змінні для нашої тестової app.

_Присутні дві опції:_

1. Однією командою підняти залежності та запустити усі тести
2. Підняти залежності, а потім окремо зупскати тести (навіть через ide)

#### Все однією командою

```bash
task back:test:func
```

#### Окремо запускати тести

1. Підянти залежності

   ```bash
   task back:test:func:up
   ```

2. Запустти тести

   ```bash
   task back:test:func:run
   ```

   АБО

   ```bash
   # після `--` вказуєте шлях до тесту
   task back:test:func:run -- ./tests/api/...
   task back:test:func:run -- ./tests/api/signup_test.go
   ```

   АБО

   **VScode та cursor (з розширенням для go) мають зруний інтерфейс, щоб запускати тести**
   **Якщо ви підняли залежнсоті на першому кроці, то зараз ви можете без проблем запускати тести з ide**

3. Опустити залежності

   ```bash
   task back:test:func:down
   ```

## Деплой

Проєкт готовий до деплою на AWS EC2 free tier.

### Швидкий старт

1. Запустіть EC2 інстанс (Ubuntu 22.04, t3.micro)
2. Підключіться через SSH та виконайте:
   ```bash
   git clone https://github.com/ukma-cs-ssdm-2025/team-circus.git
   cd team-circus
   chmod +x scripts/setup-aws.sh scripts/deploy.sh scripts/health-check.sh
   ./scripts/setup-aws.sh
   task copy:env
   nano .env  # Налаштуйте змінні оточення
   task docker:prod:deploy
   ```

### Детальна документація

- **[AWS Deployment Guide](./docs/deployment/AWS_DEPLOYMENT.md)** - Повний покроковий гайд
- **[Deployment Quick Start](./README_DEPLOYMENT.md)** - Швидкий старт

## Документація

### Обов'язкові документи для перегляду

- **[Всі артефакти вимог](./docs/requirements/)** - Повна колекція документів

- **[Документ Системного Дизайну](./docs/requirements/system-design-document.md)** - Архітектура та дизайн системи
- **[Requirements Traceability Matrix (RTM)](./docs/requirements/rtm.md)** - Матриця відстеження вимог
- **[User Stories](./docs/requirements/user-stories.md)** - Історії користувачів та вимоги
- **[Architecture Decision Records (ADR)](./docs/adr/)** - Записи архітектурних рішень
