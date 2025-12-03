import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PROTECTED_ROUTES = ["/dashboard"];

const AUTH_ROUTES = ["/auth", "/auth/register", "/auth/forgot-password"];

const PUBLIC_ROUTES = ["/", "/welcome"];

function isProtectedRoute(pathname: string): boolean {
  return PROTECTED_ROUTES.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`)
  );
}

function isAuthRoute(pathname: string): boolean {
  return AUTH_ROUTES.some((route) => pathname === route);
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/api") ||
    pathname.includes(".")
  ) {
    return NextResponse.next();
  }

  const accessToken = request.cookies.get("accessToken")?.value;
  const hasToken = Boolean(accessToken);

  if (isProtectedRoute(pathname) && !hasToken) {
    const loginUrl = new URL("/auth", request.url);
    loginUrl.searchParams.set("redirect", pathname); // Сохранение урла редиректа
    return NextResponse.redirect(loginUrl);
  }

  if (isAuthRoute(pathname) && hasToken) {
    const redirectTo = request.nextUrl.searchParams.get("redirect");
    const dashboardUrl = new URL(redirectTo || "/dashboard", request.url);
    return NextResponse.redirect(dashboardUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
};

