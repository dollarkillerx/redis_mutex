# redis_mutex
基于redis实现可靠的分布式锁

## 前言
直接使用 `SETNX` `EXPIRE`来实现分布式锁有诸多隐患,本项目就是未了解决这些隐患

### 直接使用redis作为分布式锁的隐患

#### `SETNX`, `EXPIRE` 非原子性
通讯通过TCP与redis服务进行链接 两次操作不能保证同时完成 可能造成死锁

### `EXPIRE` 设置固定超时时间是不可靠的
如果设置固定超时时间，可能当前A进程的操作还未完成，锁发生超时释放 B进程获得锁 导致数据冲突

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
