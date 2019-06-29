# curious
A go library that let you list all projects using a given import path

```go
package main

import (
  "github.com/andersfylling/curious"
  "fmt"
)

func main() {
  projects := curious.GithubSearch("github.com/andersfylling/disgord")
  fmt.Println(len(projects)) // 13
}
