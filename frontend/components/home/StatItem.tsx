type StatItemProps = {
  value: string
  label: string
}

export const StatItem = ({ value, label }: StatItemProps) => {
  return (
    <div className="flex flex-col items-center">
      <div className="text-3xl font-bold text-foreground md:text-4xl">{value}</div>
      <div className="mt-1 text-sm text-muted-foreground">{label}</div>
    </div>
  )
}

