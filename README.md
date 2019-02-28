# Recluse

This is a proof-of-concept repository that seeks to provide a simple
HTTP API w/ websockets support for a scuttlebutt data feed.

Goals:
- allow an API for FE developers to build scuttlebutt UIs without needing to understand the complexities of the lower layers
- provide a stable, purpose-driven API for building upon

Non-Goals:
- prevent people from accessing lower layers of code when/if they want them


## How to run it as a developer

1. [Install go](https://golang.org/dl/)
1. Check out this repository (*NOTE*: the default branch is `mainline`) and `cd` in to it's directory
1. `go get -v` which will download the packages onto your $GOPATH
1. Start an sbot somehow. I just run an instance of patchwork. I'm sure there are plugin things we need to figure out soon.
1. `go run *.go` will run the webservice
1. (optional) from the `demo_website` directory, type `go run *.go` which will spawn a webserver with a tiny UI for playing purposes.

## LICENSE

AGPL (if that's a problem for you, let's chat)
