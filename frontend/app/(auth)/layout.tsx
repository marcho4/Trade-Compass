import { TopNavbar } from "@/components/layout/TopNavbar"

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen bg-background">
        <TopNavbar/>
        {children}
    </div>
  )
}
