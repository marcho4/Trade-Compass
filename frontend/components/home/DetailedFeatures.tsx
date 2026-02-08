import Link from "next/link"

export const DetailedFeatures = () => {
  const features = [
    {
      emoji: "🤖",
      title: "AI-консультант, который не галлюцинирует",
      subtitle: "Спроси у AI что угодно про компанию",
      description:
        "AI читает отчетность компаний, новости, макроэкономические данные и цену акции и объясняет простым языком. Без «инвестируйте на свой страх и риск» — конкретные цифры и источники.",
      example: "Пример: «Какие риски у Газпрома на сегодняшний день?»",
      example_url: "/dashboard/GAZP",
    },
    {
      emoji: "📊",
      title: "Готовый анализ отчетов как в на Wall Street",
      subtitle: "Ключевые инсайты по компании с учетом полного контекста",
      description:
        "Больше не нужно тратить время на сбор информации и исследование рынка - нейросеть сделает всё за вас и предоставит отчет по шаблону компании с Wall Street",
      example: null,
      example_url: "",
    },
    {
      emoji: "🎯",
      title: "Сравнение компании во всём секторе друг с другом",
      subtitle: "Получайте разбор целых секторов за считанные минуты",
      description:
        "AI анализирует отчеты сразу нескольких компаний в одном секторе и приводит анализ текущего положения дел на рынке. Выделяя самую дивидендную акцию, самую перспективную и самую стабильную.",
      example: "Пример: «Проанализируй банковский сектор России»",
      example_url: "/dashboard/screener"
    },
  ]

  return (
    <section className="my-20">
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-4xl font-bold tracking-tight md:text-5xl">
          Всё, что нужно для решения — в одном месте
        </h2>
      </div>

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
                <Link
                  href={feature.example_url}
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
                </Link>
              </div>
            )}
          </div>
        ))}
      </div>
    </section>
  )
}

