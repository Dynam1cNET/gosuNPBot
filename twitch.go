package main

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rgamba/evtwebsocket"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://id.twitch.tv/oauth2/authorize",
	TokenURL: "https://id.twitch.tv/oauth2/token",
}

var (
	scopes       = []string{"chat:edit", "chat:read"}
	redirectURL  string
	oauth2Config *oauth2.Config
	twitchToken  *oauth2.Token
)

func twitch() {
	err := godotenv.Load()
	if err != nil {
		logErr("Error loading .env file")
		return
	}
	gob.Register(&oauth2.Token{})
	redirectURL = fmt.Sprintf("http://%s:%s/redirect", os.Getenv("TWITCH_REDIRECT_HOSTNAME_OR_IP"), os.Getenv("TWITCH_REDIRECT_LISTENING_PORT"))
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("TWITCH_CLIENTID"),
		ClientSecret: os.Getenv("TWITCH_SECRET"),
		Scopes:       scopes,
		Endpoint:     Endpoint,
		RedirectURL:  redirectURL,
	}

	if os.Getenv("DONT_TOUCH_TWITCH_OAUTH") != "" {
		expiretime, err := strconv.ParseInt(os.Getenv("DONT_TOUCH_UNIX_EXPIRE"), 10, 64)
		if err != nil {
			fmt.Println("Error during conversion")
			return
		}
		token := &oauth2.Token{
			AccessToken:  os.Getenv("DONT_TOUCH_TWITCH_OAUTH"),
			TokenType:    "",
			RefreshToken: os.Getenv("DONT_TOUCH_TWITCH_REFRSH"),
			Expiry:       time.Unix(expiretime, 0),
		}
		twitchToken = token

		checkKey()

		var c = evtwebsocket.Conn{
			OnConnected: func(w *evtwebsocket.Conn) {
				authToken := twitchToken.AccessToken
				chatUser := fmt.Sprintf("#%s", os.Getenv("TWITCH_STREAMER_LOGIN_NAME"))

				setAuthToken := evtwebsocket.Msg{
					Body: []byte(fmt.Sprintf("PASS oauth:%s", authToken)),
				}

				setUsername := evtwebsocket.Msg{
					Body: []byte(fmt.Sprintf("NICK %s", os.Getenv("TWITCH_BOT_LOGIN_NAME"))),
				}

				reqTags := evtwebsocket.Msg{
					Body: []byte("CAP REQ :twitch.tv/tags"),
				}

				joinChat := evtwebsocket.Msg{
					Body: []byte(fmt.Sprintf("JOIN %s", chatUser)),
				}

				err := w.Send(setAuthToken)
				err = w.Send(setUsername)
				err = w.Send(reqTags)
				err = w.Send(joinChat)
				if err != nil {
					logErr(err.Error())
				} else {
					logInfo("Connected to Twitch IRC!")
				}

			},

			OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
				go func() {
					if strings.HasPrefix(string(msg), "PING") {
						pong := evtwebsocket.Msg{
							Body: []byte(strings.Replace(string(msg), "PING", "PONG", -1)),
						}
						err := w.Send(pong)
						if err != nil {
							fmt.Printf("ERR: %s\n", err.Error())
						}
					} else {
						go messageHandler(msg, w)
					}

				}()
			},

			OnError: func(err error) {
				logErr("** ERROR **\n%s\n", err.Error())
			},

			MatchMsg: func(req, resp []byte) bool {
				return string(req) == string(resp)
			},

			Reconnect: true,
		}

		if err := c.Dial("wss://irc-ws.chat.twitch.tv:443", ""); err != nil {
			logErr(err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/redirect", HandleOAuth2Callback)
	r.HandleFunc("/", HandleLogin)
	r.HandleFunc("/done", HandleDone)

	logInfo(fmt.Sprintf("Running local Webserver on: http://%s:%s", os.Getenv("TWITCH_REDIRECT_HOSTNAME_OR_IP"), os.Getenv("TWITCH_REDIRECT_LISTENING_PORT")))
	logInfo(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("TWITCH_REDIRECT_HOSTNAME_OR_IP"), os.Getenv("TWITCH_REDIRECT_LISTENING_PORT")), r))

}

func messageHandler(msg []byte, w *evtwebsocket.Conn) {
	var parseMessage = regexp.MustCompile(`^@.*?;id=(?P<messageID>[^;]*);.*?:(?P<username>[^!]*)!.*?PRIVMSG (?P<channel>[^!]*?) :?(?P<message>.*)`)
	messageGroups := parseMessage.FindStringSubmatch(string(msg))
	if len(messageGroups) == 5 {
		messageID := messageGroups[parseMessage.SubexpIndex("messageID")]
		message := messageGroups[parseMessage.SubexpIndex("message")]
		channel := messageGroups[parseMessage.SubexpIndex("channel")]
		if strings.HasPrefix(message, fmt.Sprintf("%s%s", os.Getenv("TWITCH_COMMAND_PREFIX"), os.Getenv("TWITCH_NOW_PLAYING"))) {
			//Now Playing
			sendMessage := evtwebsocket.Msg{
				Body: []byte(fmt.Sprintf("@reply-parent-msg-id=%s PRIVMSG %s : %s", messageID, channel, fmt.Sprintf("Now Playing: %s", NowPlaying))),
			}
			err := w.Send(sendMessage)
			if err != nil {
				fmt.Println(err)
			}
		}
		if strings.HasPrefix(message, fmt.Sprintf("%s%s", os.Getenv("TWITCH_COMMAND_PREFIX"), os.Getenv("TWITCH_LAST_PLAYED"))) {
			//Last played
			sendMessage := evtwebsocket.Msg{
				Body: []byte(fmt.Sprintf("@reply-parent-msg-id=%s PRIVMSG %s : %s", messageID, channel, fmt.Sprintf("Last played maps: %s", strings.Join(LastPlayed[:], " ")))),
			}
			err := w.Send(sendMessage)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	//logInfo(fmt.Sprintf("Message from Twitch: %s", msg))

}

func HandleDone(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`<html><body><h1>You are done!</h1><p>You can close the Page now!</p></body></html>`))

	return
}

func HandleLogin(writer http.ResponseWriter, request *http.Request) {
	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		logErr("No Token :(")
	}

	state := hex.EncodeToString(tokenBytes[:])

	http.Redirect(writer, request, oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return
	}
	twitchToken = token
	var envFile, _ = godotenv.Read(".env")
	envFile["DONT_TOUCH_TWITCH_OAUTH"] = token.AccessToken
	envFile["DONT_TOUCH_TWITCH_REFRSH"] = token.RefreshToken
	envFile["DONT_TOUCH_UNIX_EXPIRE"] = fmt.Sprintf("%d", twitchToken.Expiry.Unix())
	godotenv.Write(envFile, ".env")
	logInfo("Updated .env file!")
	http.Redirect(w, r, "/done", http.StatusTemporaryRedirect)
	return
}
func checkKey() {
	type accessTokenRes struct {
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token"`
		Scope        []string `json:"scope"`
		TokenType    string   `json:"token_type"`
		Expiry       int      `json:"expires_in"`
	}
	var userdata accessTokenRes
	var now = time.Now()
	var exp = twitchToken.Expiry

	difference := exp.Sub(now)
	if difference.Hours() < 1 {
		logInfo("Twitch token needs an Update! Updating now...")
		APIUrl := "https://id.twitch.tv/oauth2/token"
		response, err := http.PostForm(APIUrl, url.Values{
			"client_id":     {oauth2Config.ClientID},
			"client_secret": {oauth2Config.ClientSecret},
			"grant_type":    {"refresh_token"},
			"refresh_token": {twitchToken.RefreshToken},
		})
		if err != nil {
			logErr("Error getting Refresh token!", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(response.Body)
		body, err := io.ReadAll(response.Body)
		err = json.Unmarshal(body, &userdata)
		twitchToken.AccessToken = userdata.AccessToken
		twitchToken.RefreshToken = userdata.RefreshToken
		twitchToken.Expiry = time.Now().Add(time.Second * time.Duration(userdata.Expiry))
		var envFile, _ = godotenv.Read(".env")
		envFile["DONT_TOUCH_TWITCH_OAUTH"] = twitchToken.AccessToken
		envFile["DONT_TOUCH_TWITCH_REFRSH"] = twitchToken.RefreshToken
		envFile["DONT_TOUCH_UNIX_EXPIRE"] = fmt.Sprintf("%d", twitchToken.Expiry.Unix())
		godotenv.Write(envFile, ".env")
		logInfo("Updated .env file!")

	}
}
