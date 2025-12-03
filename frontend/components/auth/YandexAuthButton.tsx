"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "/api";

interface YandexAuthButtonProps {
  text?: string;
  className?: string;
}


export function YandexAuthButton({
  text = "Войти через Яндекс",
  className = "w-full hover:cursor-pointer text-center",
}: YandexAuthButtonProps) {
  const [isLoading, setIsLoading] = useState(false);

  const handleYandexLogin = async () => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/auth/yandex/login`, {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to get Yandex OAuth URL");
      }

      const data = await response.json();
      
      window.location.href = data.url;
    } catch (error) {
      console.error("Yandex login error:", error);
      setIsLoading(false);
    }
  };

  return (
    <Button
      variant="outline"
      className={className}
      onClick={handleYandexLogin}
      disabled={isLoading}
      aria-label="Войти через Яндекс"
      tabIndex={0}
    >
      {isLoading ? (
        <svg
          className="mr-2 h-4 w-4 animate-spin"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      ) : (
        <svg
          className="mr-2 h-4 w-4"
          viewBox="0 0 44 44"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          aria-hidden="true"
        >
          <circle cx="22" cy="22" r="22" fill="#FC3F1D" />
          <path
            d="M24.1 32.4h3.3V11.6h-4.8c-4.9 0-7.4 2.5-7.4 6.6 0 3.1 1.4 4.9 4.1 6.8l-4.6 7.4h3.8l5.1-8.1-1.8-1.2c-2.1-1.5-3.2-2.6-3.2-5 0-2.1 1.5-3.6 4.4-3.6h1.6v17.9h-.5z"
            fill="#FFFFFF"
          />
        </svg>
      )}
      {text}
    </Button>
  );
}
