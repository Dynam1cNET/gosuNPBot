package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var NowPlaying = "User does not play anything Currently"
var LastPlayed []string

func main() {

	//Check for .env
	//If none is there we create one for the user and tell him to
	//fill it out. We also send a link to the documentation.
	if !fileExists(".env") {
		logWarn("No .env file Found!")
		templateEnv := map[string]string{
			"DONT_TOUCH_TWITCH_OAUTH":         "",
			"DONT_TOUCH_TWITCH_REFRSH":        "",
			"DONT_TOUCH_UNIX_EXPIRE":          "",
			"TWITCH_SECRET":                   "",
			"TWITCH_CLIENTID":                 "",
			"TWITCH_BOT_LOGIN_NAME":           "",
			"TWITCH_STREAMER_LOGIN_NAME":      "",
			"TWITCH_COMMAND_PREFIX":           "!",
			"TWITCH_LAST_PLAYED":              "lp",
			"TWITCH_LAST_PLAYED_HISTORY_SIZE": "5",
			"TWITCH_NOW_PLAYING":              "np",
			"TWITCH_REDIRECT_HOSTNAME_OR_IP":  "localhost",
			"TWITCH_REDIRECT_LISTENING_PORT":  "7001",
			"GOSUMEMORY_WS_IP":                "127.0.0.1",
			"GOSUMEMORY_WS_PORT":              "24050",
			"STREAMCOMPANION_WS_IP":           "localhost",
			"STREAMCOMPANION_WS_PORT":         "20727",
		}
		godotenv.Write(templateEnv, ".env")

		logInfo(".env Generated. Please fill out and Restart. See: ")
		fmt.Scanln()
		os.Exit(0)
	} else {
		logInfo("Found .env")
		godotenv.Load(".env")

	}
	logInfo("Starting local Webserver on port ")
	go twitch()
	go func() {
		for {
			streamcompanion()
			gosomemory()
		}
	}()
	fmt.Scanln()
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func logErr(s ...any) {
	dt := time.Now()
	errorCol := color.New(color.BgRed, color.FgBlack, color.Bold)
	errorCol.Printf("[ERR  | %s]", dt.Format("01-02-2006 15:04:05"))
	fmt.Print(" │ ")
	fmt.Println(s...)
}
func logInfo(s ...any) {
	dt := time.Now()
	info := color.New(color.BgCyan, color.FgBlack, color.Bold)
	info.Printf("[INFO | %s]", dt.Format("01-02-2006 15:04:05"))
	fmt.Print(" │ ")
	fmt.Println(s...)
}
func logWarn(s ...any) {
	dt := time.Now()
	warn := color.New(color.BgYellow, color.FgBlack, color.Bold)
	warn.Printf("[WARN | %s]", dt.Format("01-02-2006 15:04:05"))
	fmt.Print(" │ ")
	fmt.Println(s...)
}
