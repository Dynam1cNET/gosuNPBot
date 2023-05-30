package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func streamcompanion() {

	godotenv.Load(".env")
	type SCJson struct {
		Dl     string `json:"dl,omitempty"`
		Status int    `json:"status,omitempty"`
	}

	addr := fmt.Sprintf("%s:%s", os.Getenv("STREAMCOMPANION_WS_IP"), os.Getenv("STREAMCOMPANION_WS_PORT"))
	logInfo(fmt.Sprintf("Starting GoRoutine for StreamCompanion. Listening on: %s", addr))

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://"+addr+"/tokens", nil)
	if err != nil {
		logInfo("Error connectiong to WebSocket StreamCompanion:", err)
		return
	}
	defer conn.Close()

	jsonMessage, err := json.Marshal([...]string{"dl", "status"})
	if err != nil {
		logErr("Error Creating JSON Object:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		logErr("Error Sending JSON Payload to StreamCompanion:", err)
		return
	}

	logInfo("Connected to StreamCompanion WS:", addr)

	var lastStatus int
	var lastDL string
	arrayLimit, err := strconv.Atoi(os.Getenv("TWITCH_LAST_PLAYED_HISTORY_SIZE"))
	if err != nil {
		logErr("Error during conversion")
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logErr("StreamCompanion:", err)
			break
		}
		var data SCJson
		json.Unmarshal(message, &data)
		if data.Dl != "" {
			lastDL = data.Dl
		}
		if data.Status != -1 {
			lastStatus = data.Status
		}
		NowPlaying = lastDL
		if len(LastPlayed) > 0 {
			if LastPlayed[len(LastPlayed)-1] != lastDL && lastStatus == 2 {
				LastPlayed = append(LastPlayed, lastDL)

				if len(LastPlayed) >= arrayLimit {
					LastPlayed = LastPlayed[len(LastPlayed)-5:]
				}
				logInfo(LastPlayed)
			}
		} else {
			if lastStatus == 2 {
				LastPlayed = append(LastPlayed, lastDL)
				logInfo(LastPlayed)
			}
		}
		//logInfo("StreamCompanion:", lastDL, " MenuState:", lastStatus)
	}

	logWarn("Disconnected from StreamcompanionWS!")
}
