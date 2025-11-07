export const CTAButtons = () => {
  return (
    <div className="flex animate-fade-in-up flex-col gap-4 [animation-delay:300ms]">
      <a
        className="group inline-flex items-center justify-center rounded-full bg-primary px-10 py-5 text-lg font-bold text-primary-foreground shadow-lg transition-all hover:scale-105 hover:bg-primary/90 hover:shadow-xl focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        href="/auth/register"
        tabIndex={0}
        aria-label="Начать бесплатно"
      >
        Начать бесплатно
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
      <p className="text-sm text-muted-foreground">
        Разберем твою первую акцию за 2 минуты • Без карты
      </p>
    </div>
  )
}

