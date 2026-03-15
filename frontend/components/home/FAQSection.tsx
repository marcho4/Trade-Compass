import { ChevronDown } from "lucide-react"

export const FAQSection = () => {
  const faqs = [
    {
      question: "Какие компании покрыты?",
      answer: "10 самых крупных компаний на MOEX",
    },
    {
      question: "Как часто обновляются данные?",
      answer:
        "Метрики — каждые 6 часов. Отчетность — в день публикации.",
    },
    {
      question: "Можно ли подключить брокера?",
      answer:
        "Пока нет, но в планах на Q2 2026.",
    },
    {
      question: "AI может ошибаться?",
      answer:
        "AI помогает аггрегировать информацию, доступную на текущий день. Решение об инвестициях принимаете вы.",
    },
  ]

  return (
    <section className="my-20">
      <div className="mb-12 text-center">
        <h2 className="mb-4 text-4xl font-bold tracking-tight md:text-5xl">
          Часто задаваемые вопросы
        </h2>
      </div>

      <div className="mx-auto max-w-3xl divide-y divide-border">
        {faqs.map((faq, index) => (
          <details key={index} className="group">
            <summary className="flex cursor-pointer items-center justify-between py-5 text-left text-lg font-semibold transition-colors hover:text-primary [&::-webkit-details-marker]:hidden">
              {faq.question}
              <ChevronDown className="ml-4 h-5 w-5 flex-shrink-0 text-muted-foreground transition-transform group-open:rotate-180" aria-hidden="true" />
            </summary>
            <p className="pb-5 text-sm leading-relaxed text-muted-foreground">
              {faq.answer}
            </p>
          </details>
        ))}
      </div>
    </section>
  )
}
