import Link from "next/link";

const navigationItems: ReadonlyArray<{ label: string; href: string }> = [
  { label: "Скринер", href: "/screener" },
  { label: "Портфель", href: "/portfolio" },
];

const Header = () => {
  return (
    <header className="sticky top-6 z-50 flex w-full justify-center px-4">
      <nav
        aria-label="Главная навигация"
        className="flex w-full max-w-6xl items-center justify-between rounded-3xl border border-border bg-card/80 px-6 py-4 shadow-lg backdrop-blur-2xl backdrop-saturate-150"
      >
        <Link
          href="/"
          className="flex items-center gap-2 rounded-full bg-primary px-4 py-2 text-sm font-semibold uppercase tracking-[0.24em] text-primary-foreground shadow-sm transition-colors hover:bg-primary/90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
        >
          BullRun
        </Link>

        <div className="flex items-center gap-4">
          {navigationItems.map((item) => (
            <Link
              key={item.label}
              href={item.href}
              className="rounded-full px-4 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
            >
              {item.label}
            </Link>
          ))}
        </div>

        <Link
          href="/auth"
          className="rounded-full border border-border bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground shadow-sm transition-colors hover:bg-secondary/80 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
        >
          Войти в аккаунт
        </Link>
      </nav>
    </header>
  );
};

export default Header;

