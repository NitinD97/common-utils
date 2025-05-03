# Rate Limiter
## Example:
```go
package main

import (
	"context"
	"fmt"

	rl "github.com/jsjain/go-rate-limiter"
)

func ExampleNewLimiter() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
	InitAddress:           []string{"127.0.0.1:6379"},
  })
  if err != nil {
    panic(err)
  }
	limiter := rl.NewLimiter(client)
	res, err := limiter.Allow(ctx, "key")
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	// Output: allowed 1 remaining 9
}
```

### Setting custom rate limits for different keys and default rate limit

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alphadose/haxmap"
	rl "github.com/jsjain/go-rate-limiter"
	"github.com/redis/rueidis"
)

func NewLimiterWithCustomLimits() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		panic(err)
	}
	customLimits := haxmap.New[string, Limit]()
	customLimits.Set("key1", Limit{Burst: 50, Rate: 50, Period: time.Second})
	limiter := rl.NewLimiter(client, WithCustomLimits(customLimits), WithRateLimit(rl.PerSecond(20)))
	
	res, err := limiter.Allow(context.Background(), "key")
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	// Output: allowed 1 remaining 19

	res, err := limiter.Allow(context.Background(), "key1")
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	// Output: allowed 1 remaining 49
}

```
