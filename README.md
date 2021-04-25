# go-unique-queue
We define a queue with unique keys.
When pushing, it checks if the queue has the same key. If there is, the value will be updated, and the key is not re-queued.
When poping, it gets the latest value and executes the custom workFunc.

# required
Go 1.15
