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

- **Go** (версія 1.19+)
- **Node.js** (версія 18+)
- **Task** (taskfile.dev)

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

Достпуні дві опції:

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

## Документація

### Обов'язкові документи для перегляду

- **[Всі артефакти вимог](./docs/requirements/)** - Повна колекція документів

- **[Документ Системного Дизайну](./docs/requirements/system-design-document.md)** - Архітектура та дизайн системи
- **[Requirements Traceability Matrix (RTM)](./docs/requirements/rtm.md)** - Матриця відстеження вимог
- **[User Stories](./docs/requirements/user-stories.md)** - Історії користувачів та вимоги
- **[Architecture Decision Records (ADR)](./docs/adr/)** - Записи архітектурних рішень
