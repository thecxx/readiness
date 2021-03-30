# Readiness

> When the data we need to get is not locally, we always need to make sure that the data is ready.


## Samples

### Use readiness with default
```
import(
    "github.com/thecxx/readiness"
)

func main() {
    readiness.Register(
        "test_key",
        func(key string) (interface{}, error) {
            return "hello world", nil
        },
        2*time.Second)

    value := readiness.Get("test_key")
}

```

### Use readiness with options
```
import(
    "github.com/thecxx/readiness"
)

func main() {
    ready := readiness.New(
        readiness.WithPullFailedHandler(func(key string, err error) {
            // If there is an error when getting, we can know it
        }))
    ready.Register(
        "test_key",
        func(key string) (interface{}, error) {
            return "hello world", nil
        },
        2*time.Second)

    value := ready.Get("test_key")
}

```
