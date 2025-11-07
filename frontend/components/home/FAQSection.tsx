export const FAQSection = () => {
  const faqs = [
    {
      question: "Какие компании покрыты?",
      answer: "Все компании MOEX с публичной отчетностью. ~250 акций.",
    },
    {
      question: "Как часто обновляются данные?",
      answer:
        "Метрики — каждые 6 часов. Отчетность — в день публикации.",
    },
    {
      question: "Можно ли подключить брокера?",
      answer:
        "Пока нет, но в планах на Q2 2025. Сейчас — ручной ввод портфеля.",
    },
    {
      question: "AI может ошибаться?",
      answer:
        "Claude дает источники для каждого утверждения. Проверяй критичные решения.",
    },
  ]

  return (
    <section className="my-20">
      {/* Section Header */}
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-4xl font-bold tracking-tight md:text-5xl">
          Часто задаваемые вопросы
        </h2>
      </div>

      {/* FAQ Grid */}
      <div className="grid gap-6 md:grid-cols-2">
        {faqs.map((faq, index) => (
          <div
            key={index}
            className="rounded-2xl border border-border bg-card p-6 transition-all hover:border-primary/50"
          >
            <h3 className="mb-3 text-lg font-bold">{faq.question}</h3>
            <p className="text-sm leading-relaxed text-muted-foreground">
              {faq.answer}
            </p>
          </div>
        ))}
      </div>
    </section>
  )
}


