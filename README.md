# loadbalancer

a loadbalancer for routing hashes to different servers

TODO:

- implement round robin
- implement based on resources (later on)
- might need to change hardcoded byte 1024
- split the code into different files
- add error handling for doing a get command with a -1 port
- just like aws, a feature would be to add a new server if the current server load is full, potentially scale down if the load is low as well but this will be only if either all the keys are expirable keys or if the server is not being used at all or if data is not stored in the server memory. This will also require a solution for getting ports that are not the hosts ports and allocate a range of say maybe 10,000 ports and docker can create new servers when needed so a possible solution is to use a different network instead of 127.0.0.1

-1 {"command": "get", "key": "test"}
-1 {"command": "set", "key": "test", "value": "testvalue"}

TODO: receive requests and route then to a port via round robin
and then return the port back to the client

the client will need to keep track of the port
the client should have load balancing off by default

FIXME:

- need to allow multiple inputs
