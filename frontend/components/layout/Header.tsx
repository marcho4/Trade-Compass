import Link from "next/link"

export const Header = ({ showLogin = true }: { showLogin?: boolean }) => {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-7 h-[52px] border-b border-border bg-card">
      <div className="flex items-center gap-4">
        <Link
          href="/"
          className="font-mono text-[13px] font-bold text-primary tracking-[2px] no-underline"
        >
          Trade Compass
        </Link>

        <div className="flex items-center gap-2.5 font-mono text-[11px] text-muted-foreground tracking-[0.5px]">
          <span className="text-muted-foreground/50">&gt;</span>
          <span className="text-foreground">ГЛАВНАЯ</span>
        </div>
      </div>

      {showLogin && (
        <Link
          href="/auth"
          className="inline-flex items-center gap-[7px] font-mono text-[11px] font-semibold tracking-[0.8px] text-muted-foreground hover:text-primary no-underline px-[11px] py-[5px] border border-border hover:border-primary rounded-[2px] bg-transparent hover:bg-primary/10 transition-all duration-[120ms]"
        >
          ВОЙТИ
        </Link>
      )}
    </header>
  )
}
