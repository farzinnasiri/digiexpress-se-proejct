# DLocator

---

## Overview

We have a courier locator microservice. this microservice is responsible for returning the couriers present within a
given haversine radius of a certain location. in order to achieve such functionality it simply listens on a stream of gps points,
storing the last location of each courier. and when the query rpc is called, it queries the current locations of the couriers and
returns the result. to get a better understanding see [locator.proto](./api/proto/v1/locator.proto)

This service is not implemented yet, this repo only contains a bunch of boilerplate code, we want you to implement it.
See the incomplete code in [handler.go](./internal/app/dlocator/handler.go).

---

## Requirements

- Correctness
  - the query rpc must return correct results (please write tests for it)
  - some locations are expired. e.g. imagine a gps point for 1 day ago, this gps point is clearly useless.
  - we want the excess gps points to be cleared up and have TTLs.
  - so they wouldn't show up in the future query results.
- Minimal Latency
    - there is a benchmark inside this repository. you can run it with `make bench`
    - we have a solution that has <1ms (sub millisecond) latency per operation for our query api.
    - you are encouraged to achieve this latency. but it is ok to not achieve it.
    - try your best.
- Easy Deployment
  - we need to easily deploy this service and its dependencies via docker(compose).
      - define `up` target inside Makefile so that `make up` deploys everything.
      - define `down` target inside Makefile so that `make down` stops everything.
- Code quality
    - we need a simple, easy to understand and clean code.
    - some code style issues can be detected using linters.
    - we have put a makefile target for you to lint your code. run it with `make lint`
- Recovery
  - when this microservice is restarted for any reason (CI/CD pipelines, configuration changes, machine failure etc.) we want the courier data to remain.
  - meaning that the data cannot be lost

---

## Constraints

- do not change the boilerplate code
- you can only spawn at max one goroutine in each rpc
- you can use and deploy any library, or service (e.g. mysql redis osm etc.)
- this microservice will only be deployed in one replica. you can avoid the distributed challenges for simplicity.
- after you're done implementing the service, write a paragraph or two inside this README.md and elaborate how to scale your service and deploy it on multiple replicas.

---

## How to run

- run `make dev-tools` to install linters, formatters, etc.
- run `make protoc` to generate protobuf client and server
- run `make dependencies` to download and install module dependencies
- run `make wire` to generate dependency injection code
- run `make help` to see other make targets descriptions.

---

## Contact us

- if you have any questions of any sort, please contact the assigned engineer via his/her phone.
