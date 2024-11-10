# CSE 138 - Assignment 3

### Team Members:

Madison Li, Margaret Heathcote, & Wyatt Hawes

## Overview & Instructions

Create a replicated, fault-tolerant, and causally consistent key-value store.

...

#

...

```
go run http_server.go helper_funcs.go key_value_ops.go
set VIEW="localhost:8090,localhost:8091"
```
...

```
docker build -t app .
docker run --rm -p 8090:8090 -e=VIEW=localhost:8090,localhost:8091 app
```

## Mechanism Description

...(how our system tracks causal dependencies)...

...(how our system detects when a replica goes down)...

## Files

- `Dockerfile`: a Dockerfile describing how to create a container to build and run our code.
- `README.md`: a description of the project, instructions, team contributions, and citations.
- `helper_funcs.go`: ...
- `http_server.go`: ...
- `key_value_ops.go`: ...
- `test_assignment3.py`: provided python test script.

## Team Contributions

**Wyatt** - ...

**Maggie** - ...

**Madison** - ...

## Acknowledgements

We did not consult anyone on this assignment.

## Citations

- ...