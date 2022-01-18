## Intro

Bot for MUD game written in Go, primarily for educational purposes. Heavily underimplemented and at alpha stage.

Currently based on RMUD (rmud.ru).

## Run

1. `./mudbot <local host:port> <mud host:port>`
2. Connect via any MUD client to local host port.

Bot will work in "assisted mode" — it will proxy MUD output to client as-is (uncompressed), but will use it to perform internal logic.

## Bot logic

Currently only parsing "score" output and updating current status.

## Bot commands

TBD
