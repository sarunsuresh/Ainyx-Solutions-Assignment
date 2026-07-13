package websocket

import (
    "github.com/gofiber/contrib/websocket"
    "github.com/gofiber/fiber/v2"
)


func UpgradeMiddleware() fiber.Handler {
    return websocket.New(nil)
}

func Handler(hub *Hub) fiber.Handler {
    return websocket.New(func(conn *websocket.Conn) {
        hub.Register(conn)
        defer hub.Unregister(conn)

        for {
            _, _, err := conn.ReadMessage()
            if err != nil {
                break  
            }
        }
    })
}