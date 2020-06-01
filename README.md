Challenge:  Third Term Partial
==============================

This is the third term challenge-based exam for the Distributed Computing Class. This is the beginning of 3 deliverables of your final challenge.
On this challenge you're developing the first part of your final project that will be a parallel image processing system.

A strong recomendation is that you develop your solution the most simple, readable, scalable and plugable as possible. In the following challenges you will
be integrating more services, so a well-defined design and implementation will make it easier to integrate new modules into your distributed application.

Distributed and Parallel Image Processing
-----------------------------------------

![architecture](architecture.png)

### Second Phase for the Final Challage
This is going to be the second phase of design and implementation.
On this phase you are adding 3 new components that will start making more sense as a distributed system
- Controller
- Scheduler
- Worker

Your project will be divided on packages with very descriptive names where each system's component will be implemented.
Below you can see the details of each package and requirements for this partial:

- `api/`
  - From now, all request must be token-based authenticated
  - **Endpoint:** `/status` - Overall system status and logged user details. Also, print workers details (name, status and usage percentage)
  - **Endpoint:** `/status/<worker>` - Per worker details:  (name, tags, status and usage percentage)
  - **Endpoint:** `workloads/test` - This endpoint will trigger an initial end-to-end test from the `api` to a running `worker`

- `controller/`
  - Basic overall system and per node data store (it can be in-memory or a key-value datastore)
  - Request pre-validation before sending to `scheduler`
  - Controller will create a message-passing server for its interaction with workers

- `scheduler/`
  - Basic workloads scheduling that will be based on node's tags and # of running workloads
  - Scheduler is calling workers through RPC

- `worker/`
  - Standalone component with initial `test` RPC function.
  - Worker's command line will be as follows:
    - `./worker --controller <host>:<port> --node-name <node_name> --tags <tag1>,<tag2>...`


**Documentation**
- A detailed arquitecture document will be required for this initial phase in the [architecture.md](architecture.md) file. Diagrams and charts can be included on this document.
- A detailed user guide must be written in the [user-guide.md](user-guide.md) file. This document explains how to install, configure and use your system.


Test Cases (from console)
-------------------------
- [Project's First Phase Test Cases](../second-partial/#test-cases-from-console)

- **Node Status**
```
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status/<worker>
{
	"Worker": "Worker-name",
	"Tags": "tag1,tag2,tag3",
	"Status": "Running",
	"Usage": "50%"
}
```

- **Execute Workload**
```
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads/test
{
	"Workload": "test",
	"Job ID": "1",
	"Status": "Scheduling",
	"Result: "Done in Worker: <worker_name>"
}
```

"Game" Rules
------------

- This is 2-person team challenge, keep the focus on you work.
- You're free to use the internet for coding references.
- Any attempt of plagiarism will not be tolerated.


General Submission Instructions
-------------------------------
1. Make sure your local repository is in sync with the origin remote repository before anything.
2. Commit and Push your code to your personal repository (fork) and branch (first-partial).

3. Once you're done, follow common lab's sumission process. More details at: [Classify API](../../classify.md)
```
GITHUB_USER=<your_github_account> make submit

# Example:
GITHUB_USER=obedmr make submit
```

Grading Policy
--------------

The grading policy is quite simple, most falls in the test cases. Below the percentages table:

| Concept                                | %    |
|----------------------------------------|------|
| Code Style best practices              | 20%  |
| Test Cases (one for each API endpoint) | 60%  |
| Program meets with all requirements    | 20%  |
| TOTAL                                  | 100% |

Handy links
-----------
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Postman](https://www.postman.com/)
- [Video: Basics of Using Postman](https://youtu.be/t5n07Ybz7yI)
- [Advanced REST client for Chrome browser](https://chrome.google.com/webstore/detail/advanced-rest-client/hgmloofddffdnphfgcellkdfbfbjeloo?hl=es-419)