export const ObjectionsSection = () => {
  const objections = [
    {
      question: "А почему не Finam/Smart-Lab бесплатно?",
      answer: "За экономию времени. AI прочитает отчет на 200 страниц за вас.",
      benefits: [
        "AI-саммари отчета на 200 страниц за 30 секунд",
        "Не нужно переключаться между сайтами — всё в одном месте",
        "Удобный интерфейс вместо таблиц из 90-х",
      ],
    },
    {
      question: "Зачем платить, если есть ChatGPT?",
      answer:
        "ChatGPT не знает, где искать отчетность MOEX, и может ошибаться из-за отстутствия нужного контекста",
      benefits: [
        "Прямой доступ к актуальным отчетам",
        "Метрики уже рассчитаны и обновляются автоматически",
        "Сравнение с индустрией из нашей базы",
        "AI ассистентом с полным пониманием компании"
      ],
    },
  ]

  return (
    <section className="my-20">
      <div className="grid gap-12 md:grid-cols-2">
        {objections.map((objection, index) => (
          <div
            key={index}
            className="rounded-2xl border border-border bg-card p-8"
          >
            <h3 className="mb-4 text-2xl font-bold">{objection.question}</h3>

            <p className="mb-6 text-lg font-semibold text-primary">
              {objection.answer}
            </p>

            {objection.benefits && (
              <div className="space-y-3">
                {objection.benefits.map((benefit, idx) => (
                  <div key={idx} className="flex items-start gap-2">
                    <svg
                      className="mt-0.5 h-5 w-5 flex-shrink-0 text-green-500"
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
                    <p className="text-sm leading-relaxed">{benefit}</p>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>
    </section>
  )
}


