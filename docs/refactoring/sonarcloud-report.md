# SonarCloud Результати

## До: 
![](sonar-before1.png)

### Трошки почистивши False Positive's, маємо наступне:
![](sonar-before2.png)

## Після:
![](sonar-after.png)

## Обрані метрики:
| Метрика           | До    | Після |
| ----------------- | ----- | ----- |
| Reliability       | 2, C  |   -   |
| Maintainability   | 33, A |   -   |
| Security Hotspots | 7     |   -   |
| Duplications      | 6.9%  |   -   |

### Застосовувались наступні патерни: 
- **Replace Magic Number** – прибрали магічні числа, що стосувались секретного ключа та термінів придатності JSON Web токенів

## Регресія:
![](image.png)