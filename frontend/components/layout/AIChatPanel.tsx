"use client"

import * as React from "react"
import { Send, Sparkles, Trash2, MessageSquare } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { cn } from "@/lib/utils"
import { PromptExample, defaultPrompts } from "@/lib/ai-prompts"

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  timestamp: Date
}

interface AIChatPanelProps {
  promptExamples?: PromptExample[]
}

export const AIChatPanel = ({ 
  promptExamples = defaultPrompts
}: AIChatPanelProps) => {
  const [messages, setMessages] = React.useState<Message[]>([
    {
      id: "1",
      role: "assistant",
      content: "Привет! Я твой AI-ассистент. Чем могу помочь в анализе портфеля?",
      timestamp: new Date(),
    },
  ])
  const [input, setInput] = React.useState("")
  const [isLoading, setIsLoading] = React.useState(false)
  const messagesEndRef = React.useRef<HTMLDivElement>(null)
  const messagesContainerRef = React.useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    if (messagesContainerRef.current) {
      messagesContainerRef.current.scrollTo({
        top: messagesContainerRef.current.scrollHeight,
        behavior: "smooth"
      })
    }
  }

  React.useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSendMessage = async () => {
    if (!input.trim() || isLoading) return

    const userMessage: Message = {
      id: Date.now().toString(),
      role: "user",
      content: input.trim(),
      timestamp: new Date(),
    }

    setMessages((prev) => [...prev, userMessage])
    setInput("")
    setIsLoading(true)

    // Имитация ответа AI (здесь будет интеграция с реальным API)
    setTimeout(() => {
      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: `Получил ваш запрос: "${userMessage.content}". Это демо-ответ. Скоро здесь будет настоящая интеграция с AI.`,
        timestamp: new Date(),
      }
      setMessages((prev) => [...prev, aiMessage])
      setIsLoading(false)
    }, 1000)
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const handleClearChat = () => {
    setMessages([
      {
        id: "1",
        role: "assistant",
        content: "Привет! Я твой AI-ассистент. Чем могу помочь в анализе портфеля?",
        timestamp: new Date(),
      },
    ])
  }

  const handlePromptClick = (prompt: string) => {
    setInput(prompt)
  }

  return (
    <div className="flex h-full flex-col bg-card border border-border rounded-2xl shadow-lg overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-border px-4 py-3">
        <div className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-full bg-gradient-to-br from-primary via-[hsl(var(--chart-1))] to-[hsl(var(--chart-3))] flex items-center justify-center shadow-md">
            <Sparkles className="h-4 w-4 text-primary-foreground" />
          </div>
          <h2 className="font-semibold text-base">AI Ассистент</h2>
        </div>
        <Button
          variant="ghost"
          size="icon"
          onClick={handleClearChat}
          className="h-8 w-8"
          aria-label="Очистить чат"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>

      {/* Messages Area */}
      <div ref={messagesContainerRef} className="flex-1 overflow-y-auto px-4 py-4 space-y-4 min-h-0">
        {messages.length === 0 ? (
          <div className="flex h-full items-center justify-center text-center">
            <div className="space-y-3">
              <MessageSquare className="mx-auto h-12 w-12 text-muted-foreground" />
              <p className="text-sm text-muted-foreground">
                Начните диалог с AI-ассистентом
              </p>
            </div>
          </div>
        ) : (
          messages.map((message) => (
            <div
              key={message.id}
              className={cn(
                "flex gap-3",
                message.role === "user" ? "justify-end" : "justify-start"
              )}
            >
              {message.role === "assistant" && (
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                  <Sparkles className="h-4 w-4 text-primary" />
                </div>
              )}
              <Card
                className={cn(
                  "max-w-[75%] px-3 py-2",
                  message.role === "user"
                    ? "bg-primary text-primary-foreground"
                    : "bg-muted"
                )}
              >
                <p className="text-sm leading-relaxed whitespace-pre-wrap">
                  {message.content}
                </p>
                <p
                  className={cn(
                    "mt-1.5 text-xs",
                    message.role === "user"
                      ? "text-primary-foreground/70"
                      : "text-muted-foreground"
                  )}
                >
                  {message.timestamp.toLocaleTimeString("ru-RU", {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </p>
              </Card>
              {message.role === "user" && (
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary">
                  <span className="text-xs font-medium text-primary-foreground">
                    Вы
                  </span>
                </div>
              )}
            </div>
          ))
        )}
        {isLoading && (
          <div className="flex gap-3">
            <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
              <Sparkles className="h-4 w-4 text-primary animate-pulse" />
            </div>
            <Card className="bg-muted px-3 py-2">
              <div className="flex gap-1">
                <span className="h-2 w-2 animate-bounce rounded-full bg-primary [animation-delay:-0.3s]" />
                <span className="h-2 w-2 animate-bounce rounded-full bg-primary [animation-delay:-0.15s]" />
                <span className="h-2 w-2 animate-bounce rounded-full bg-primary" />
              </div>
            </Card>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      <Separator />

      {/* Prompt Examples */}
      <div className="px-4 py-3 space-y-2 bg-muted/20">
        <p className="text-xs font-medium text-muted-foreground mb-2">
          Примеры вопросов:
        </p>
        <div className="grid grid-cols-2 gap-2">
          {promptExamples.map((example) => (
            <button
              key={example.id}
              onClick={() => handlePromptClick(example.prompt)}
              className={cn(
                "p-3 rounded-lg border border-border bg-background",
                "hover:bg-accent hover:border-primary/40 transition-all duration-200",
                "text-left text-xs font-medium",
                "focus:outline-none focus:ring-2 focus:ring-primary/50"
              )}
              aria-label={`Использовать пример: ${example.title}`}
            >
              <div className="flex items-start gap-2">
                <Sparkles className="h-3 w-3 text-primary/70 shrink-0 mt-0.5" />
                <span className="line-clamp-2">{example.title}</span>
              </div>
            </button>
          ))}
        </div>
      </div>

      <Separator />

      {/* Input Area */}
      <div className="px-4 py-3">
        <div className="flex gap-2">
          <Input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Задайте вопрос..."
            className="flex-1"
            disabled={isLoading}
            aria-label="Введите сообщение"
          />
          <Button
            onClick={handleSendMessage}
            disabled={!input.trim() || isLoading}
            size="icon"
            aria-label="Отправить сообщение"
          >
            <Send className="h-4 w-4" />
          </Button>
        </div>
        <p className="mt-2 text-xs text-muted-foreground">
          AI может совершать ошибки
        </p>
      </div>
    </div>
  )
}

