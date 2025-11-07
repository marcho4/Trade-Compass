import { ReactNode } from "react"

type FeatureCardProps = {
  icon: ReactNode
  title: string
  description: string
  iconBgColor: string
  iconColor: string
}

export const FeatureCard = ({
  icon,
  title,
  description,
  iconBgColor,
  iconColor,
}: FeatureCardProps) => {
  return (
    <div className="group rounded-3xl border border-border bg-card p-8 shadow-lg backdrop-blur-xl transition-all hover:scale-[1.02] hover:border-primary/30 hover:shadow-xl">
      <div
        className={`mb-4 inline-flex h-12 w-12 items-center justify-center rounded-2xl ${iconBgColor} ${iconColor}`}
      >
        {icon}
      </div>
      <h2 className="text-2xl font-semibold text-card-foreground">{title}</h2>
      <p className="mt-3 leading-relaxed text-muted-foreground">{description}</p>
    </div>
  )
}

