package main

import (
	"log"
	"os"
	"os/signal"

	gowebsocket "github.com/sacOO7/GoWebsocket"
)

func main() {
	/*resp, err := http.Get("https://golang.org")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	io.Copy(os.Stdout, resp.Body)
	*/

	/*resp, err := http.Post(
		"https://postman-echo.com/post",
		"application/json",
		strings.NewReader(`{"foo1":"bar1","foo2":"bar2"}`),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	*/

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("ws://echo.websocket.org/")
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Соединение установлено")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Соединение не установлено", err)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Получен пинг" + data)
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("Закрытие соединения с сервером")
			socket.Close()
			return
		}
	}
}
