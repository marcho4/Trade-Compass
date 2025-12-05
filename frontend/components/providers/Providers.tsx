"use client";

import { AuthProvider } from "@/contexts/AuthContext";
import { CookieBanner } from "@/components/CookieBanner";
import { type ReactNode } from "react";

interface ProvidersProps {
  children: ReactNode;
}

export function Providers({ children }: ProvidersProps) {
  return (
    <AuthProvider>
      {children}
      <CookieBanner />
    </AuthProvider>
  );
}

