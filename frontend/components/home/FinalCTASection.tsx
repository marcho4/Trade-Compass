export const FinalCTASection = () => {
  return (
    <section className="my-20">
      <div className="overflow-hidden rounded-3xl border border-primary/20 bg-gradient-to-br from-primary/10 via-primary/5 to-background p-12 text-center shadow-2xl md:p-16">
        {/* Heading */}
        <h2 className="mb-6 text-4xl font-bold tracking-tight md:text-5xl">
          Перестань гуглить метрики по 10 сайтам
        </h2>

        {/* Description */}
        <p className="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground md:text-xl">
          Bull Run собирает данные, AI объясняет их смысл, ты принимаешь
          решения. Как должно быть.
        </p>

        {/* CTA Button */}
        <a
          href="/auth/register"
          className="group inline-flex items-center justify-center rounded-full bg-primary px-10 py-5 text-lg font-bold text-primary-foreground shadow-lg transition-all hover:scale-105 hover:bg-primary/90 hover:shadow-xl focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          tabIndex={0}
          aria-label="Проанализировать первую акцию бесплатно"
        >
          Проанализировать первую акцию бесплатно
          <svg
            className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1"
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

        {/* Subtext */}
        <p className="mt-6 text-sm text-muted-foreground">
          2 минуты на анализ • Без карты • Отмена в 1 клик
        </p>
      </div>
    </section>
  )
}


