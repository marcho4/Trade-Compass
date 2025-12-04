export const SocialProofSection = () => {
  return (
    <section className="my-20">
      <div className="mx-auto max-w-4xl rounded-2xl border border-border bg-gradient-to-br from-card to-card/50 p-8 md:p-12">
        <div className="mb-6">
          <span className="inline-flex items-center rounded-full bg-primary/10 px-4 py-1.5 text-sm font-semibold text-primary">
            Founder's Note
          </span>
        </div>

        <blockquote className="mb-8 space-y-4 text-lg leading-relaxed text-foreground">
          <p>
            Я потерял 380к рублей в 2023-м, купив «дешевые» акции по P/E.
            Проблема была не в метрике, а в том, что я не понимал контекст.
          </p>
          <p>
            Эта платформа — то, чего мне не хватало тогда: данные + их
            объяснение в одном месте. Без необходимости быть аналитиком.
          </p>
        </blockquote>

        <div className="flex items-center gap-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gradient-to-br from-primary to-primary/60 text-xl font-bold text-primary-foreground">
            M
          </div>
          <div>
            <p className="font-semibold">Марк</p>
            <p className="text-sm text-muted-foreground">Founder, Trade Compass</p>
          </div>
        </div>
      </div>
    </section>
  )
}


