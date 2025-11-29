import GridBackground from "@/components/ui/grid-background"
import { CTAButtons } from "./CTAButtons"
import { StatsSection } from "./StatsSection"

export const HeroSection = () => {
  return (
    <GridBackground className="min-h-[700px] overflow-hidden rounded-[2.5rem] border border-border shadow-xl md:h-[700px]">
      {/* Content */}
      <div className="relative z-20 flex h-full flex-col items-center justify-center px-4 py-12 text-center sm:px-6 md:px-12 md:py-16">
        <div className="flex max-w-5xl flex-col items-center gap-6 md:gap-8">
          {/* Heading */}
          <h1 className="max-w-4xl animate-fade-in-up bg-linear-to-br from-foreground to-foreground/70 bg-clip-text text-5xl font-bold leading-tight text-transparent [animation-delay:100ms] sm:text-4xl md:text-6xl lg:text-7xl">
            Анализируй российские акции как профи на Уолл-стрит
          </h1>

          {/* Description */}
          <p className="max-w-3xl animate-fade-in-up text-lg leading-relaxed text-muted-foreground [animation-delay:200ms] md:text-xl">
            AI разберет отчетность любой компании MOEX за 2 минуты и покажет, что видят институционалы. Без финансового образования.
          </p>

          {/* CTA Buttons */}
          <CTAButtons />

          {/* Stats */}
          <StatsSection />
        </div>
      </div>
    </GridBackground>
  )
}

