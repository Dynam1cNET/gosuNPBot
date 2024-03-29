# gosuNPBot
Simple osu NP bot for your Twitch chat

This twitch bot uses either StreamCompanion or gosumemory to get a !np and !lp command in your Twitch chat!

Setup is as simple as it gets. 
Download executable, fill out Twitch bot data and start. 


### Setup: 
- Download precompiled executable from [releases](https://github.com/Dynam1cNET/gosuNPBot/releases/).
- Start executable once. It will generate a .env file. 
- The .env is mostly pre-filled. If you never changes settings on gosu / streamcompanion then you are good to go!
- Go to https://dev.twitch.tv/console/ and register a application.
- If you use the pre-filled .env then the "OAuth Redirect URLs" should be "http://localhost:7001/redirect".
- After you proceed you see a list of applications (or just one). Edit it and click on new Secret. 
- Copy Client-ID in the .env "TWITCH_CLIENTID".
- Copy Client-Secret in the .env "TWITCH_SECRET".
- Edit "TWITCH_BOT_LOGIN_NAME" and "TWITCH_LOGIN_NAME". If you are running the bot as your normal Twitch user, you need to put your name in both. 
- Go to http://localhost:7001/ - This should redirect you to Twitch. Accept your own application and you should get redirected to a "Success page". Restart the bot once and you are done. You should now see ```Connected to Twitch IRC!``` in the console.
- Edit the other stuff as you please. See .env below for help.

You should be ready to go. Run the tool and it should connect to the Twitch IRC and StreamCompanion / gosumemory. 

### .env

The .env should look something like this:
```
DONT_TOUCH_TWITCH_OAUTH=""
DONT_TOUCH_TWITCH_REFRSH=""
DONT_TOUCH_UNIX_EXPIRE=""
GOSUMEMORY_WS_IP="127.0.0.1"
GOSUMEMORY_WS_PORT=24050
STREAMCOMPANION_WS_IP="localhost"
STREAMCOMPANION_WS_PORT=20727
TWITCH_BOT_LOGIN_NAME=""
TWITCH_CLIENTID=""
TWITCH_COMMAND_PREFIX="\!"
TWITCH_LAST_PLAYED="lp"
TWITCH_LAST_PLAYED_HISTORY_SIZE=5
TWITCH_NOW_PLAYING="np"
TWITCH_REDIRECT_HOSTNAME_OR_IP="localhost"
TWITCH_REDIRECT_LISTENING_PORT=7001
TWITCH_SECRET=""
TWITCH_STREAMER_LOGIN_NAME=""
```
- DONT_TOUCH_TWITCH_OAUTH: This will get filled out by the bot itself. 
- DONT_TOUCH_TWITCH_REFRSH: This will get filled out by the bot itself. 
- DONT_TOUCH_UNIX_EXPIRE: This will get filled out by the bot itself. 
- GOSUMEMORY_WS_IP: This should be the ip gosu is running on. Normally 127.0.0.1
- GOSUMEMORY_WS_PORT: This is the port gosu is listening in. Default is 24050
- STREAMCOMPANION_WS_IP: This should be the ip SC is running on. Normally localhost (can also be 127.0.0.1)
- STREAMCOMPANION_WS_PORT: This is the port SC is listening in. Default is 20727
- TWITCH_BOT_LOGIN_NAME: If you use another Twitch account for your bot, you need to type its name here. Else just your Twitch login name
- TWITCH_CLIENTID: The Client-ID for your application. You can get it from https://dev.twitch.tv/console/
- TWITCH_COMMAND_PREFIX: The prefix for your command. Default is "!" so a command could be !np. If you choose & a command could be &np
- TWITCH_LAST_PLAYED: The command for last played. Last played gives out a list of the x last played maps. 
- TWITCH_LAST_PLAYED_HISTORY_SIZE: How many it should give out. Five should be enough. I recommend to not go above. 
- TWITCH_NOW_PLAYING: The command for "Now playing". This also works in the menu. It just gives out the song you are currently listening to or playing right now. 
- TWITCH_REDIRECT_HOSTNAME_OR_IP: This is for your Twitch application. You should only change it, if you know what you are doing. :D
- TWITCH_REDIRECT_LISTENING_PORT: The port for your Twitch application. same as above. ^^
- TWITCH_SECRET: The Client-Secret for your application. You can get it from https://dev.twitch.tv/console/
- TWITCH_STREAMER_LOGIN_NAME: This should be just your Twitch name. That channel the bot will go to and listen for Twitch commands. 
