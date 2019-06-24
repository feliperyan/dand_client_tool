# D&D DM Remote Tool 🧙
## Client app to connect to the D&D [websocket server](https://github.com/feliperyan/dand_server_tool) to play streaming MP3s and receive messages in a classic UI of such poor taste it's kinda cool.

> 🚨This project was a fun way for me to get more comfortable with Golang and most of all to grok Go routines and channels. It's very much in the _make it work_ phase 🍝 of the "make it work -> make it right -> make it fast" concept.

### First deploy the server somewhere, [here is the repo](https://github.com/feliperyan/dand_server_tool)

### Regular players
Start up the executable, enter the servername and hit enter, then enter your player name and hit enter.

### Dungeon Master
Start the tool from the terminal and add the flags:

```--mode dm --addr <server_address>```

You have acess to the following commands, which must be started with a ```/``` backslash:
- ```/audio <file.mp3>``` streams an mp3 file to all players and start playing it.
- ```/list``` lists the name of all players connected.
- ```/whisper <playerName>``` sends a message just to that specific player.
- ```<text>``` broadcasts a message to all players.

> This little repo would not be possible without the work by Hajime Hoshi and the [Ebiten](https://github.com/hajimehoshi/ebiten) Go package.

