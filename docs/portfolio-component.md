# Компонент состава портфеля

## Обзор

Реализован функционал отображения портфелей инвестора с детальным составом позиций.

## Структура

### 1. Компоненты

#### `PortfolioComposition` 
**Путь:** `frontend/components/portfolio/PortfolioComposition.tsx`

Компонент для отображения детального состава портфеля с позициями.

**Возможности:**
- Отображение списка всех позиций в портфеле
- Расчет текущей стоимости каждой позиции
- Расчет прибыли/убытка по каждой позиции (в рублях и процентах)
- Отображение доли каждой позиции в портфеле
- Визуализация веса позиции с помощью прогресс-бара
- Сортировка позиций по весу в портфеле (от большего к меньшему)
- Отображение информации о секторе, количестве лотов, средней и текущей цене

**Props:**
```typescript
interface PortfolioCompositionProps {
  positions: Position[]
  totalValue: number
}
```

### 2. Страницы

#### Главная страница портфелей
**Путь:** `frontend/app/(main)/dashboard/portfolio/page.tsx`

Отображает список всех портфелей пользователя в виде карточек.

**Особенности:**
- ✅ Убраны графики производительности
- ✅ Убран селектор риска
- ✅ Убраны настройки целей
- ✅ Добавлена информационная подсказка о необходимости выбрать портфель
- ✅ Клик по карточке портфеля ведет на детальную страницу

#### Детальная страница портфеля
**Путь:** `frontend/app/(main)/dashboard/portfolio/[id]/page.tsx`

Отображает детальную информацию о выбранном портфеле.

**Секции:**
1. **Навигация** - кнопка возврата к списку портфелей
2. **Заголовок** - название и описание портфеля
3. **Состав портфеля** - компонент `PortfolioComposition`
4. **График производительности** - сравнение с индексами
5. **Настройки риска** - выбор уровня риска (conservative/moderate/aggressive)
6. **Цели портфеля** - установка и отслеживание целей
7. **Информация о стратегии** - полезные советы

### 3. Типы данных

#### `Position`
**Путь:** `frontend/types/portfolio.ts`

Структура данных для позиции в портфеле, соответствующая таблице `position` из БД:

```typescript
interface Position {
  id: string
  portfolioId: string
  companyId: string
  companyTicker: string
  companyName: string
  quantity: number // количество лотов
  avgPrice: number // средняя цена покупки
  currentPrice: number // текущая цена
  lastBuyDate: Date
  sector?: string
  createdAt?: Date
  updatedAt?: Date
}
```

#### `Portfolio`
**Путь:** `frontend/types/portfolio.ts`

Структура данных для портфеля, соответствующая таблице `portfolio` из БД:

```typescript
interface Portfolio {
  id: string
  name: string
  userId: string
  description?: string
  value: number
  createdAt: Date
  updatedAt?: Date
  profitPercent?: number
  profitAmount?: number
  rating?: number
  positions?: Position[]
}
```

### 4. Демо-данные

**Путь:** `frontend/lib/mock-data.ts`

Централизованный файл с демо-данными для разработки и тестирования:

- `mockPositions` - массив демо-позиций для разных портфелей
- `mockPortfolios` - объект с информацией о портфелях
- `generatePerformanceData()` - функция генерации данных графика
- `getPositionsByPortfolioId(portfolioId)` - получить позиции портфеля
- `getPortfolioById(portfolioId)` - получить информацию о портфеле

## Соответствие схеме БД

### Таблица `portfolio`
```dbml
Table portfolio {
  id serial pk
  name varchar(100)
  user_id int [ref: > users.id, not null]
  description text
  created_at timestamptz [default: `now()`]
  updated_at timestamptz
}
```

### Таблица `position`
```dbml
Table position {
  id serial pk
  portfolio_id int [ref: > portfolio.id]
  company_id int [ref: - companies.id]
  avg_price decimal(10, 2)
  last_buy_date timestamptz
  quantity int // кол-во лотов
  created_at timestamptz [default: `now()`]
  updated_at timestamptz
}
```

## Расчеты

### Стоимость позиции
```typescript
positionValue = currentPrice * quantity
```

### Прибыль/убыток позиции
```typescript
profit = (currentPrice - avgPrice) * quantity
profitPercent = (profit / (avgPrice * quantity)) * 100
```

### Доля позиции в портфеле
```typescript
weight = (positionValue / totalPortfolioValue) * 100
```

## UI/UX особенности

1. **Адаптивность**: Компонент адаптирован для мобильных устройств и десктопов
2. **Цветовая индикация**: Зеленый для прибыли, красный для убытка
3. **Иконки**: TrendingUp/TrendingDown для визуального отображения тренда
4. **Сортировка**: Автоматическая сортировка позиций по весу в портфеле
5. **Прогресс-бары**: Визуальное отображение доли каждой позиции

## Следующие шаги

1. **Интеграция с API**: Подключить реальные данные из бэкенда
2. **Редактирование позиций**: Добавить возможность редактирования/удаления позиций
3. **Добавление новых позиций**: Форма для добавления компаний в портфель
4. **Ребалансировка**: Рекомендации по оптимизации состава портфеля
5. **История транзакций**: Отображение истории покупок/продаж
6. **Фильтрация и поиск**: Возможность фильтровать позиции по секторам

