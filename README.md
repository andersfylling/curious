# curious
A go library that let you list all projects using a given import path. 


> Please register a github token (no permissions needed) to the environment variable "CURIOUS_GITHUB_TOKEN".

```go
package main

import (
  "github.com/andersfylling/curious"
  "fmt"
)

func main() {
  projects, err := curious.GithubSearch("github.com/andersfylling/disgord")
  if err != nil {
  	panic(err)
  }
  
  fmt.Println(len(projects)) // 13
}
```
