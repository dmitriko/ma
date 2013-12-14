ma
==

Manage set of workers via websockets protocol.
First stage of development is to get an
application for tailing/greping logs on remote hosts.

Tasks:  
* server accepts websocket connection from remote hosts
* server accepts a websocket connection only from allowed IP
* server processes YAML, JSON, TOML format for config files
* server has simple api to add/remove allowed IP
