# CSE 138 - Assignment 3

### Team Members:

Madison Li, Margaret Heathcote, & Wyatt Hawes

## Overview & Instructions

Created a replicated, fault-tolerant, and causally consistent key-value store.

#

- Build container image and tag it `asg3img`:

```
$ docker build -t asg3img .
```

- Create subnet called asg3net with IP range 10.10.0.0/16:

```
$ docker network create --subnet=10.10.0.0/16 asg3net
```

- Run each replica in the subnet:

```
$ docker run --rm -p 8082:8090 --net=asg3net --ip=10.10.0.2 --name=alice -e=SOCKET_ADDRESS=10.10.0.2:8090 -e=VIEW=10.10.0.2:8090,10.10.0.3:8090,10.10.0.4:8090 asg3img

$ docker run --rm -p 8083:8090 --net=asg3net --ip=10.10.0.3 --name=bob -e=SOCKET_ADDRESS=10.10.0.3:8090 -e=VIEW=10.10.0.2:8090,10.10.0.3:8090,10.10.0.4:8090 asg3img

$ docker run --rm -p 8084:8090 --net=asg3net --ip=10.10.0.4 --name=carol -e=SOCKET_ADDRESS=10.10.0.4:8090 -e=VIEW=10.10.0.2:8090,10.10.0.3:8090,10.10.0.4:8090 asg3img
```

## Mechanism Description

#### How our system tracks causal dependencies:

- We keep track of causal dependencies by giving each key in the key-value-pair map a version number. Every time a request is sent or an update is received, this version is sent along with it. This allows each server to know if it is up to date and which events on a key happened after another, since if the version is greater, then it is more recent. We decided to tie-break conflicts by simply comparing the IP address of the two servers in question, the lower IP address always takes priority. This way both servers can tie-break without talking to each other.

#### How our system detects when a replica goes down:

- Our system detects when a replica goes down because each replica periodically (every 2.5 seconds) sends a request to all other replicas via a version of the VIEW that contains every replica ever added to the view. Essentially, each replica periodically notifies all other replicas to add the sender to its VIEW (where, generally, it should already be inside). If a server does not respond to this request, it's considered "down" and deleted from the dynamic version of the view (which does remove downed replicas). When a server comes back online, its VIEW request that it periodically sends will be received by the replicas, adding it back into their VIEW.

## Files

- `Dockerfile`: a Dockerfile describing how to create a container to build and run our code.
- `README.md`: a description of the project, instructions, team contributions, and citations.
- `helper_funcs.go`: code for helper functions.
- `http_server.go`: main code that starts an http server and handles each endpoint.
- `key_value_ops.go`: functions for key-value operations (get, put, etc).
- `view_ops.go`: functions for view operations (add, delete, etc).
- `test_assignment3.py`: provided python test script.

## Team Contributions

We met multiple times as a group to discuss the assignment and collaborate on code.

**Wyatt** - Collaborated using VSCode live share; worked on key-value store functions

**Maggie** - Collaborated using VSCode live share; worked on view functions

**Madison** - Collaborated using VSCode live share; worked on version and validating

## Acknowledgements

We did not consult anyone on this assignment.

## Citations

- Writing coroutines/periodic actions in Go: https://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals

- GO http package documentation: https://pkg.go.dev/net/http

- How to get started with go: https://go.dev/doc/tutorial/getting-started
