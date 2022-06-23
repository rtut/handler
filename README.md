## Example main

```
package main

import (
    ch "github.com/rtut/handler"
    "log"
    "net/http"
)

func main() {
    http.Handle("/", &ch.Handler{})
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## About errors and tests

The quantity and quality of code coverage with tests depends on many factors:
* requirements (the task does not say about the percentage of code coverage by tests, do I need to open critical parts of the code with tests or not)
* capabilities (time limit on develop, difficult cases, different type of tests)
* etc .....

May be necessary test on timeout (and logic in code on goroutine timeout).
Part of code where i have problem with http.Head, necessary test it based on expected behavior:
* ignore problem with http.Head request (i use it and didn't handle this behavior in test)
* interrupt all program

My cover 91.2%.

```
go tool cover -html=my_coverage.out
```

Test on race condition, did not reveal a data race.

```
go tool -race .
```
