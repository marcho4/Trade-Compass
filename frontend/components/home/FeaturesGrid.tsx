import { FeatureCard } from "./FeatureCard"
import { ComputerIcon } from "./icons/ComputerIcon"
import { FilterIcon } from "./icons/FilterIcon"
import { ChartIcon } from "./icons/ChartIcon"

export const FeaturesGrid = () => {
  const features = [
    {
      icon: <ComputerIcon />,
      title: "AI-Аналитика",
      description:
        "30 секунд на автоматический отчет компании с рекомендациями и метриками.",
      iconBgColor: "bg-gradient-to-br from-slate-500/10 to-slate-600/10",
      iconColor: "text-slate-600 dark:text-slate-400",
    },
    {
      icon: <FilterIcon />,
      title: "Умный скринер",
      description:
        "Фильтры по отраслям, мультипликаторам и динамике с обновлением раз в минуту.",
      iconBgColor: "bg-gradient-to-br from-slate-600/10 to-slate-700/10",
      iconColor: "text-slate-700 dark:text-slate-300",
    },
    {
      icon: <ChartIcon />,
      title: "Баланс портфеля",
      description:
        "Трекинг позиций, автоматические уведомления и подсказки по ребалансировке.",
      iconBgColor: "bg-gradient-to-br from-slate-700/10 to-slate-800/10",
      iconColor: "text-slate-800 dark:text-slate-200",
    },
  ]

  return (
    <div className="grid gap-6 md:grid-cols-3">
      {features.map((feature) => (
        <FeatureCard
          key={feature.title}
          icon={feature.icon}
          title={feature.title}
          description={feature.description}
          iconBgColor={feature.iconBgColor}
          iconColor={feature.iconColor}
        />
      ))}
    </div>
  )
}
