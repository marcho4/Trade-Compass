package application

import "context"

type ChatManager interface {
	CreateChat(ctx context.Context, userId string, systemPrompt string) (string, error) // метод, который будет создавать чат и возвращать ID чата
	GetUserChats(ctx context.Context, userId string) ([]string, error)                  // метод, который будет возвращать список чатов для пользователя
	GetChatMessages(ctx context.Context, chatId string) ([]string, error)               // метод, который будет возвращать список сообщений для чата
	DeleteChat(ctx context.Context, chatId string) error                                // метод, который будет удалять чат
	GetPreparedMessages() ([]string, error)                                             // метод, который будет давать варианты для сообщений
	TransformMessage(ctx context.Context, message string) (string, error)               // для конвертации сообщения в ембеддинг для векторного поиска
	SearchContext(ctx context.Context, message string) (string, error)                  // для поиска контекста в векторном хранилище
	SendMessage(ctx context.Context, chatId string, message string) (string, error)     // метод, который будет отправлять сообщение в чат и возвращать ID сообщения
	GetResponseStream(ctx context.Context, chatId string) (string, error)               // метод, который будет возвращать поток ответов для чата
}
