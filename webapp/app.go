package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func main() {
    router := gin.Default()
    router.LoadHTMLFiles("index.html")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    router.GET("/ws", func(c *gin.Context) {
        handleWebSocket(c.Writer, c.Request)
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    go consumeKafka()
    go handleMessages()

    // Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()
    clients[conn] = true

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            delete(clients, conn)
            break
        }
    }
}

func consumeKafka() {
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "172.18.0.3:31551", 
        "group.id":          "main-consumer-group",
        "auto.offset.reset": "earliest",
    })

    if err != nil {
        fmt.Printf("Failed to create consumer: %s", err)
        os.Exit(1)
    }

    topic := "test"
    err = c.SubscribeTopics([]string{topic}, nil)
    if err != nil {
        fmt.Printf("Failed to subscribe to topics: %s", err)
        os.Exit(1)
    }

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    run := true
    for run {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev, err := c.ReadMessage(100 * time.Millisecond)
            if err != nil {
                continue
            }
            message := fmt.Sprintf("Consumed event from topic %s: key = %-10s value = %s",
                *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
            broadcast <- message
        }
    }

    c.Close()
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, []byte(msg))
            if err != nil {
                log.Printf("WebSocket error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}
