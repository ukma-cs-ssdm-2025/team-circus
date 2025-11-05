# Звіт результатів SonarCloud

## До рефакторингу
![SonarCloud metrics before initial cleanup](sonar-before1.png)

### Трохи почистивши False Positives, маємо наступне:
![SonarCloud metrics after false positives cleanup](sonar-before2.png)

## Після рефакторингу
![SonarCloud metrics after refactoring](sonar-after.png)

## Обрані метрики

| Метрика             | До     | Після |
|---------------------|--------|:-----:|
| Reliability         | 2, C   | 0, A  |
| Maintainability     | 33, A  | 30, A |

### Застосовувались наступні патерни: 
- **Replace Magic Number** – прибрано магічні числа, що стосувалися секретного ключа та термінів придатності JSON Web Token
- **Simplify Conditional** - перетворення складної логіки на простішу

## Перевірка регресії
![Check for regresson](regression-check.png)
*З нового pull request-у після рефакторингу.