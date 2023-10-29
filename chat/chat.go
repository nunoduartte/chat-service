package chat

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"net/http"
)

type Chat struct {
	RedisClient *redis.Client
}

func (c *Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roteador := http.NewServeMux()
	roteador.Handle("/", http.HandlerFunc(c.template))
	roteador.Handle("/send", http.HandlerFunc(c.sendMessage))
	roteador.Handle("/chat", http.HandlerFunc(c.webSocket))
	roteador.ServeHTTP(w, r)
}

func (c *Chat) template(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "chat.html")
}

func (c *Chat) webSocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer conn.Close()

	roomName := r.URL.Query().Get("room")

	pubsub := c.RedisClient.Subscribe(r.Context(), roomName)
	defer pubsub.Close()

	redisChannel := pubsub.Channel()

	for msg := range redisChannel {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			fmt.Println("Erro ao enviar mensagem WebSocket:", err)
			return
		}
	}
}

func (c *Chat) sendMessage(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("room")
	message := r.URL.Query().Get("message")

	err := c.RedisClient.Publish(r.Context(), roomName, message).Err()
	if err != nil {
		http.Error(w, "Erro ao publicar mensagem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
