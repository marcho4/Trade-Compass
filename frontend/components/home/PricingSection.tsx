import Link from "next/link"

export const PricingSection = () => {
  const plans = [
    {
      name: "FREE",
      price: "Бесплатно",
      description: "Для пробы и разовых проверок",
      features: [
        "3 компании бесплатно",
        "Базовые метрики всех компаний",
        "Разбор отчетов за последний год",
      ],
      cta: "Начать бесплатно",
      ctaLink: "/auth/register",
      popular: false,
    },
    {
      name: "Pro",
      price: "19990₽/квартал",
      originalPrice: "39990₽",
      description: "Для активных инвесторов. Помимо базовых метрик по компании, вы получаете следующее:",
      features: [
        "Безлимитный доступ к нейроассистенту от Google по всем отчетам на платформе - 5000₽ ценность",
        "Безлимитный доступ к анализу отчетов от топовых LLM моделей - 200 000₽ ценность",
        "Ведение портфеля акций с анализом от нейроассистента - 6000₽ ценность",
        "Персональная поддержка в телеграме от основателя со средним временем ответа не больше 10 минут",
        "Безусловная гарантия возврата денег, если не устроил сервис"
      ],
      cta: "Стать Pro",
      ctaLink: "/auth/register",
      popular: true,
    }
  ]

  return (
    <section className="my-20">
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-4xl font-bold tracking-tight md:text-5xl">
          Начни бесплатно, плати когда нужно больше
        </h2>
      </div>

      <div className="grid gap-8 md:grid-cols-2">
        {plans.map((plan, index) => (
          <div
            key={index}
            className={`relative overflow-hidden rounded-2xl border p-8 transition-all hover:shadow-lg ${
              plan.popular
                ? "border-primary bg-primary/5 shadow-md"
                : "border-border bg-card"
            }`}
          >
            {plan.popular && (
              <div className="absolute right-4 top-4">
                <span className="inline-flex items-center rounded-full bg-primary px-3 py-1 text-xs font-bold text-primary-foreground">
                  ПОПУЛЯРНО
                </span>
              </div>
            )}

            <h3 className="mb-2 text-sm font-bold uppercase tracking-wider text-muted-foreground">
              {plan.name}
            </h3>

            <div className="mb-4">
              {plan.originalPrice && (
                <p className="text-sm text-muted-foreground line-through">
                  {plan.originalPrice}
                </p>
              )}
              <p className="text-4xl font-bold">{plan.price}</p>
            </div>

            <p className="mb-6 text-sm text-muted-foreground">
              {plan.description}
            </p>

            <ul className="mb-8 space-y-3">
              {plan.features.map((feature, idx) => (
                <li key={idx} className="flex items-start gap-2">
                  <svg
                    className="mt-0.5 h-5 w-5 flex-shrink-0 text-primary"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth={2}
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M5 13l4 4L19 7"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                  <span className="text-sm">{feature}</span>
                </li>
              ))}
            </ul>

            <Link
              href={plan.ctaLink}
              className={`block w-full rounded-full py-3 text-center text-sm font-semibold transition-all hover:scale-105 ${
                plan.popular
                  ? "bg-primary text-primary-foreground shadow-lg hover:bg-primary/90"
                  : "border border-border bg-secondary text-secondary-foreground hover:bg-secondary/80"
              }`}
              tabIndex={0}
              aria-label={plan.cta}
            >
              {plan.cta}
            </Link>
          </div>
        ))}
      </div>
    </section>
  )
}


