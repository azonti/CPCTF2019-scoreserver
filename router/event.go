package router

import (
	"fmt"
	"net/http"

	"github.com/trevex/golem"
)

var (
	Ws                   *golem.Router
	Room                 = golem.NewRoomManager()
	openProblemEventChan chan openProblemEvent
	sendFlagEventChan    chan sendFlagEvent
)

type openProblemEvent struct {
	EventName string `json:"eventName"`
	UserID    string `json:"userID"`
	ProblemID string `json:"problemID"`
}

type sendFlagEvent struct {
	EventName string `json:"eventName"`
	UserID    string `json:"userID"`
	Username  string `json:"username"`
	ProblemID string `json:"problemID"`
	IsSolved  bool   `json:"isSolved"`
}

func connClose(conn *golem.Connection) {
	Room.LeaveAll(conn)
}

func connOpen(conn *golem.Connection, req *http.Request) {
	Room.Join("event", conn)
}

func SetupWs() error {
	Ws = golem.NewRouter()
	if err := Ws.OnClose(connClose); err != nil {
		panic(err)
	}
	if err := Ws.OnConnect(connOpen); err != nil {
		panic(err)
	}
	Ws.On("po", func(conn *golem.Connection) {
		fmt.Println("popopop")
	})

	openProblemEventChan = make(chan openProblemEvent)
	sendFlagEventChan = make(chan sendFlagEvent)

	go func() {
		for {
			select {
			case event := <-openProblemEventChan:
				Room.Emit("event", "", event)

			case event := <-sendFlagEventChan:
				Room.Emit("event", "", event)
			}
		}
	}()

	return nil
}
