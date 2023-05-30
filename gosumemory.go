package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func gosomemory() {
	godotenv.Load(".env")
	type GosoMemoryJson struct {
		Menu struct {
			State int `json:"state,omitempty"`
			Bm    struct {
				ID  int `json:"id,omitempty"`
				Set int `json:"set,omitempty"`
			} `json:"bm,omitempty"`
		} `json:"menu,omitempty"`
	}

	addr := fmt.Sprintf("%s:%s", os.Getenv("GOSUMEMORY_WS_IP"), os.Getenv("GOSUMEMORY_WS_PORT"))
	logInfo(fmt.Sprintf("Starting GoRoutine for gosuMemory. Listening on: %s", addr))
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://"+addr+"/ws", nil)
	if err != nil {
		logInfo("Error connectiong to WebSocket gosuMemory:", err)
		return
	}
	defer conn.Close()

	logInfo("Connected to gosuMemory WS:", addr)
	arrayLimit, err := strconv.Atoi(os.Getenv("TWITCH_LAST_PLAYED_HISTORY_SIZE"))
	if err != nil {
		logErr("Error during conversion")
		return
	}
	var tmpData GosoMemoryJson
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logErr("gosuMemory:", err)
			break
		}
		var data GosoMemoryJson
		json.Unmarshal(message, &data)
		if tmpData != data {
			NowPlaying = fmt.Sprintf("http://osu.ppy.sh/b/%d", data.Menu.Bm.ID)
			if len(LastPlayed) > 0 {
				if LastPlayed[len(LastPlayed)-1] != fmt.Sprintf("http://osu.ppy.sh/b/%d", data.Menu.Bm.ID) && data.Menu.State == 2 {
					LastPlayed = append(LastPlayed, fmt.Sprintf("http://osu.ppy.sh/b/%d", data.Menu.Bm.ID))
					//logInfo("gosuMemoryReader:", fmt.Sprintf("http://osu.ppy.sh/b/%d", data.Menu.Bm.ID), " MenuState:", data.Menu.State)

					if len(LastPlayed) >= arrayLimit {
						LastPlayed = LastPlayed[len(LastPlayed)-5:]
					}
					logInfo(LastPlayed)
				}
			} else {
				if data.Menu.State == 2 {
					LastPlayed = append(LastPlayed, fmt.Sprintf("http://osu.ppy.sh/b/%d", data.Menu.Bm.ID))
					logInfo(LastPlayed)
				}
			}

			tmpData = data

		}

	}

	logWarn("Disconnected from gosuMemoryWS!")
}
