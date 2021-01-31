# ActiveWorkflow Agent Example in Go

This project implements a simple example agent for [ActiveWorkflow](https://github.com/automaticmode/active_workflow).
The agent is implemented in Go and uses the [remote agent API](https://docs.activeworkflow.org/remote-agent-api).

This agent doesn't do anything particularly useful but demonstrates some of the
features of ActiveWorkflow's Remote Agent API. The agent simply records the the
number of `check` and `receive` calls it sees in memory and returns messages
with the current count.

## Quick Start

This agent is intended to be used as part of [ActiveWorkflow](https://github.com/automaticmode/active_workflow).
To get started with ActiveWorkflow please see the [ActiveWorkflow documentation](https://docs.activeworkflow.org/).

You will need a working Go toolchain installed. There are no external
dependencies. The agent is built with just the Go standard library.

Please make sure you run this agent *before* starting ActiveWorkflow.  The
agent can be started without a compile step, like this:

```sh
go run agent.go
```

Alternatively, you can compile it first and then run it:

```sh
go build -o agent .
./agent
```

Note the URL of the agent's server (usually `http://127.0.0.1:5000/`), and set
it as an environment variable for ActiveWorkflow:

```sh
export REMOTE_AGENT_URL="http://127.0.0.1:5000/"
```
Now you can start ActiveWorkflow. You should be able to create instances of
this agent (named "Go Test Counter Agent"). Run it and send messages to it.

If using Docker to run ActiveWorkflow, you'll need to use the `-e` parameter to
`docker run` to pass `REMOTE_AGENT_URL` through to ActiveWorkflow. The address
in the URL will also have to be updated to match where the agent is running
(`127.0.0.1` is unlikely to be correct).

Please note that this project is just a minimal example. Consider using a
proper project structure when developing your own ActiveWorkflow agents.
