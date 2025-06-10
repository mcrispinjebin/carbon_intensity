# carbon_intensity

The Carbon Intensity API is a RESTful service that provides the forecast period with the lowest carbon intensity for a specified duration.
It supports retrieving both continuous and non-continuous time slots.
---

[Go](https://go.dev/)  
[Docker](https://www.docker.com/)

---

**Contents**

1. [Setup](#setup)
1. [Assumptions](#assumptions)
1. [Quality](#quality)
1. [Future Scope](#future-scope)

---

### Setup ###

1. Install Golang and ensure Go project can be run in the system.
1. Install docker to run the application.
1. Clone the repo in local
1. Use the command `make build` to build the docker image.
1. Use the command `make run` to run the docker image at port 3000.
1. curl the service with the desired input to get the response. e.g., `http://localhost:3000/slots?duration=120&continuous=false`
1. Use the command `make test` to run the unit tests.

---


### Assumptions ###
1. Altered the response to accommodate both the slots and calculated average.
1. The average calculation is kept as integer for simplicity.
1. If partial duration is provided, the service will find the partial duration only at the end of the window.

---

### Quality ###

Unit Test cases are available in `processor/processor_test` file.

---

### Future Scope ###

1. Partial duration calculation to be improve by checking it on each slot on the window.
1. Add a linting tool, for static check.
1. Documentation of APIs in OpenAPI(Swagger)
1. Improving unit test coverage.

---

