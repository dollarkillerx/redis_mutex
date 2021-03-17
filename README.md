# redis_mutex
Implementing reliable distributed locks based on redis

## Preface
Directly using `SETNX` `EXPIRE` to implement distributed locks has many pitfalls, and this project is to solve these pitfalls

### The pitfalls of using redis directly as a distributed lock

#### `SETNX`, `EXPIRE` are non-atomic
Communication is linked to the redis service over TCP Two operations cannot be guaranteed to complete simultaneously May cause deadlocks

### `EXPIRE` Setting a fixed timeout is unreliable
If a fixed timeout is set, it is possible that the current operation of process A has not yet completed and the lock is released on timeout process B gets the lock resulting in data conflict

### use 
```go
	log.SetFlags(log.Llongfile | log.LstdFlags)

	mutex, err := New(SetUri("127.0.0.1", 6379))
	if err != nil {
		log.Fatalln(err)
	}

	look, err := mutex.Look()
	if err != nil {
		log.Fatalln(err)
	}

	if !look {
		log.Fatalln("Lock Failure look")
	}

	time.Sleep(time.Second * 10)

	unLook, err := mutex.UnLook()
	if err != nil {
		log.Fatalln(err)
	}

	if !unLook {
		log.Fatalln("Lock Failure unLook")
	}

	time.Sleep(time.Second * 10)
```