# go-whatsup-client

[![Code Quality & Tests](https://github.com/sprak3000/go-whatsup-client/actions/workflows/quality-and-tests.yml/badge.svg)](https://github.com/sprak3000/go-whatsup-client/actions/workflows/quality-and-tests.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/f378f20d8587cd169e69/maintainability)](https://codeclimate.com/github/sprak3000/go-whatsup-client/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/f378f20d8587cd169e69/test_coverage)](https://codeclimate.com/github/sprak3000/go-whatsup-client/test_coverage)

This client allows you to retrieve the status page data for various services. It currently supports the following
status page formats:

- statuspage.io JSON responses (Example: [reddit Status](https://www.redditstatus.com/api/v2/status.json))
- [Slack API 2.0 JSON responses](https://api.slack.com/docs/slack-status#v2_0_0__current-status-api)

## Usage

### Retrieving statuspage.io style status pages

The URL for a Status page using the [statuspage.io](https://www.atlassian.com/software/statuspage) format typically
follows the pattern `https://<domain>/api/v2/status.json`. You provide a service name and a URL for the status page to
`whatsup.StatuspageIoService()`. Using GitHub as an example, we use `github` for the service name and the URL
`https://www.githubstatus.com/api/v2/status.json`.

```go
package main

import (
	"fmt"
	"os"

	"github.com/sprak3000/go-whatsup-client/whatsup"
)

func main() {
	v, err := whatsup.StatuspageIoService("github", "https://www.githubstatus.com/api/v2/status.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%#v\n", v)
}
```

If an error occurs, the error message contains the `github` service name as part of the message. The return value is a
`statuspageio.Response` structure.

### Retrieving Slack's status page

As there is only one Slack status page, you simply call `whatsup.Slack()`.

```go
package main

import (
	"fmt"
	"os"

	"github.com/sprak3000/go-whatsup-client/whatsup"
)

func main() {
	v, err := whatsup.Slack()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%#v\n", v)
}
```

If an error occurs, the error message contains `slack` as the service name as part of the message. The return value is a
`slack.Response` structure.

## Contributing

You want to contribute to the project? Welcome!

Since this is an open source project, we would love to have your feedback! If you are interested, we would also love to
have your help! Whether helpful examples to add to the docs, or FAQ entries, everything helps. Read our guide on
[contributing](docs/contributing.md), then [set up the tooling](docs/development.md) necessary.
