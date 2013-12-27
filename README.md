ma
==

Manage set of workers via websockets protocol.
First stage of development is to get an
application for tailing/greping logs on remote hosts.

Tasks:  
* server accepts websocket connection from remote hosts
* server accepts a websocket connection only from allowed IP
* server processes YAML, JSON format for config files
* server has simple api to add/remove cluster config


Cluster config has name, description and list of hosts.
Eash host is a config with remote_ip, roles, title fields.

