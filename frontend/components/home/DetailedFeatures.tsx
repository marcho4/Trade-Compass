import Link from "next/link"
import { Bot, BarChart3, Target, ArrowRight } from "lucide-react"
import type { LucideIcon } from "lucide-react"

type Feature = {
  icon: LucideIcon
  title: string
  subtitle: string
  description: string
  example: string | null
  example_url: string
}

export const DetailedFeatures = () => {
  const features: Feature[] = [
    {
      icon: Bot,
      title: "AI-консультант, который не галлюцинирует",
      subtitle: "Спроси у AI что угодно про компанию",
      description:
        "AI читает отчетность компаний, новости, макроэкономические данные и цену акции и объясняет простым языком. Без «инвестируйте на свой страх и риск» — конкретные цифры и источники.",
      example: "Пример: «Какие риски у Газпрома на сегодняшний день?»",
      example_url: "/dashboard/GAZP",
    },
    {
      icon: BarChart3,
      title: "Готовый анализ отчетов как на Wall Street",
      subtitle: "Ключевые инсайты по компании с учетом полного контекста",
      description:
        "Больше не нужно тратить время на сбор информации и исследование рынка - нейросеть сделает всё за вас и предоставит отчет по шаблону компании с Wall Street",
      example: null,
      example_url: "",
    },
    {
      icon: Target,
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
            <div className="mb-4">
              <feature.icon className="h-12 w-12 text-primary" />
            </div>

            <h3 className="mb-2 text-2xl font-bold">{feature.title}</h3>

            <p className="mb-4 text-base font-semibold text-primary">
              {feature.subtitle}
            </p>

            <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
              {feature.description}
            </p>

            {feature.example && (
              <div className="mt-4">
                <Link
                  href={feature.example_url}
                  className="inline-flex items-center text-sm font-medium text-primary transition-colors hover:text-primary/80"
                  aria-label={feature.example}
                >
                  {feature.example}
                  <ArrowRight className="ml-1 h-4 w-4" aria-hidden="true" />
                </Link>
              </div>
            )}
          </div>
        ))}
      </div>
    </section>
  )
}
