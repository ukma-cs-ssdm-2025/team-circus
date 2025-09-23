# Markdown Circus Docs

[![CI Test](https://github.com/ukma-cs-ssdm-2025/team-circus/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/ukma-cs-ssdm-2025/team-circus/actions/workflows/ci.yml)

Наш проект MCD (Markdown Circus Docs) - це інструмент спільного редагування документів у форматі Markdown. Детальний опис проєкту [тут](./Project-Description.md)

**Чому ми обрали саме Markdown?**

- Легко та швидко редагується
- Має уніфікований і стандартизований синтаксис
- Підтримується більшістю IDE та редакторів
- Знайомий і зрозумілий кожному члену команди розробки

## Учасники команди

- Артур Ключка
  - Repo Maintainer, Requirements Lead
  - github: [velosypedno](https://github.com/velosypedno)
- Владислав Панько
  - CI Maintainer, Quality Lead
  - github: [xncff](https://github.com/xncff)
- Олексій Костик
  - Documentation Lead
  - github: [kitmelancholic](https://github.com/kitmelancholic)
- Ярослав Шах
  - Issue Tracker Lead, Traceability Lead
  - github: [YaroslavShakh](https://github.com/YaroslavShakh)

Командний статут знаходиться [тут](./TeamChapter.md).

## Деталі розробки

- Стратегію бранчування - [github flow](https://docs.github.com/en/get-started/using-github/github-flow)
- Стратегія неймінгу комітів - [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)
- Мова програмування обрана для розробки серверної частини - [Go](https://go.dev/)
- Мова програмування обрана для розробки клієнтської частини - [Typescript](https://www.typescriptlang.org/) ([React](https://react.dev/))

## Структура репозиторію

```
.github/workflows/ci.yml
docs/
| requirements/          - артифакти вимог
loom/                    - відео-презентації
.gitignore
Project-Description.md
README.md
TeamChapter.md
```
