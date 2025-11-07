export const ObjectionsSection = () => {
  const objections = [
    {
      question: "А почему не Finam/Smart-Lab бесплатно?",
      answer: "Finam показывает цифры. Мы показываем, что они значат.",
      details: [
        {
          source: "Finam",
          text: "P/E = 5.2",
        },
        {
          source: "Bull Run",
          text: "P/E = 5.2 — в 2 раза ниже среднего по энергетике. Либо рынок недооценивает (аргументы), либо есть скрытые риски (аргументы). Аналоги в секторе: ...",
        },
      ],
      additional: "У нас AI-анализ документов, которого нет ни у кого в России.",
    },
    {
      question: "Зачем платить, если есть ChatGPT?",
      answer:
        "ChatGPT не знает, где искать отчетность MOEX, и выдает устаревшие данные.",
      benefits: [
        "Прямой доступ к актуальным отчетам",
        "Метрики уже рассчитаны и обновляются автоматически",
        "Сравнение с индустрией из нашей базы",
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
            {/* Question */}
            <h3 className="mb-4 text-2xl font-bold">{objection.question}</h3>

            {/* Answer */}
            <p className="mb-6 text-lg font-semibold text-primary">
              {objection.answer}
            </p>

            {/* Details (for Finam comparison) */}
            {objection.details && (
              <div className="mb-6 space-y-4">
                {objection.details.map((detail, idx) => (
                  <div
                    key={idx}
                    className="rounded-lg border border-border bg-background/50 p-4"
                  >
                    <p className="mb-2 text-sm font-semibold text-muted-foreground">
                      {detail.source}:
                    </p>
                    <p className="text-sm leading-relaxed">{detail.text}</p>
                  </div>
                ))}
              </div>
            )}

            {/* Benefits (for ChatGPT comparison) */}
            {objection.benefits && (
              <div className="mb-6 space-y-3">
                <p className="text-sm font-semibold text-muted-foreground">
                  Bull Run:
                </p>
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

            {/* Additional info */}
            {objection.additional && (
              <p className="text-sm font-medium text-primary">
                + {objection.additional}
              </p>
            )}
          </div>
        ))}
      </div>
    </section>
  )
}


