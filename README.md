# loadbalancer

a loadbalancer for routing hashes to different servers

TODO:

- implement round robin
- implement based on resources (later on)
- might need to change hardcoded byte 1024
- split the code into different files
- add error handling for doing a get command with a -1 port

-1 {"command": "get", "key": "test"}
-1 {"command": "set", "key": "test", "value": "testvalue"}

TODO: receive requests and route then to a port via round robin
and then return the port back to the client

the client will need to keep track of the port
the client should have load balancing off by default

FIXME:

- need to allow multiple inputs
