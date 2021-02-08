# GoTwitchRouter
This is a twitch cmd router. 
The router receives messages via twitch IRC and routes (messages which begin with !) them to apps.
Apps register via GRPC to the router with the command, a accessLevel and the help message.
The router will then return a GRPC stream that receives messages.

The cmds: part, join and usernotice are not reachable from chat these cmds relay the events that are equally named, see twitch IRC documentation.
# WIP
A lot of work has to be done to make this router stable, so use at your own risk.

# Requirements
## run
- postgresql
## develop
- protoc
- protoc-go
- protoc-go-grpc
- sqlc
# Install

execute the schemas in the ./asset/sqlc/schema folder in your database.

look at the config.yml from where you run it and set db credentials, and twitch infos.

the twitch oauth token should be stored in the TWITCH env variable.

# TODO
- AccessLevel restrictions (almost done)
- ...(create an issue if you want/need more)



