export const DetailedFeatures = () => {
  const features = [
    {
      emoji: "🤖",
      title: "AI, который не галлюцинирует",
      subtitle: "Спроси у AI что угодно про компанию",
      description:
        "Claude 3.5 читает отчетность, ищет в интернете и объясняет простым языком. Без «инвестируйте на свой страх и риск» — конкретные цифры и источники.",
      example: "Пример: «Почему Газпром упал на 15%?»",
    },
    {
      emoji: "📊",
      title: "Визуализация как в мировых платформах",
      subtitle: "Сравни с индустрией одним взглядом",
      description:
        "Snowflake-оценка (как Simply Wall St), графики метрик против конкурентов, автообновление каждые 6 часов. Не таблицы Excel, а понятные визуалы.",
      example: null,
    },
    {
      emoji: "🎯",
      title: "Идеи по ребалансировке",
      subtitle: "Получай идеи по улучшению портфеля",
      description:
        "AI анализирует твой портфель и предлагает, где риски, где переконцентрация, что докупить для диверсификации. Как личный аналитик за 1990₽ вместо 50к/месяц.",
      example: null,
    },
  ]

  return (
    <section className="my-20">
      {/* Section Header */}
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-4xl font-bold tracking-tight md:text-5xl">
          Всё, что нужно для решения — в одном месте
        </h2>
      </div>

      {/* Features Grid */}
      <div className="grid gap-8 md:grid-cols-3">
        {features.map((feature, index) => (
          <div
            key={index}
            className="group relative overflow-hidden rounded-2xl border border-border bg-card p-8 transition-all hover:border-primary/50 hover:shadow-lg"
          >
            {/* Icon */}
            <div className="mb-4 text-5xl">{feature.emoji}</div>

            {/* Title */}
            <h3 className="mb-2 text-2xl font-bold">{feature.title}</h3>

            {/* Subtitle */}
            <p className="mb-4 text-base font-semibold text-primary">
              {feature.subtitle}
            </p>

            {/* Description */}
            <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
              {feature.description}
            </p>

            {/* Example Link */}
            {feature.example && (
              <div className="mt-4">
                <a
                  href="#"
                  className="inline-flex items-center text-sm font-medium text-primary transition-colors hover:text-primary/80"
                  tabIndex={0}
                  aria-label={feature.example}
                >
                  {feature.example}
                  <svg
                    className="ml-1 h-4 w-4"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth={2}
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M13 7l5 5m0 0l-5 5m5-5H6"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                </a>
              </div>
            )}
          </div>
        ))}
      </div>
    </section>
  )
}

