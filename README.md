# introduction
We define a queue with unique keys.

When pushing, it checks if the queue has the same key. If there is, the value will be updated, and the key is not re-queued.

When poping, it gets the latest value and executes the custom workFunc defined by yourself.


# required
Go 1.15

# example

```go

func workFunc(obj interface{}) error {
  fmt.Println(obj.(string))
}

func main() {
  // parameters:
  // capacity: the queue capacity
  // worker: the number of worker
  // workFunc: when pop, the worker will execute the work function
  // onSuccess: when workFunc execute succeed, the onSuccess func will be executed
  // onFailed: when workFunc execute failed, the onFailed func will be executed
  // stopCh：the stop signal

  q = queue.NewQueue(100, 5, workFunc, nil, nil, nil)
}

```

