package redis_mutex

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// 设置锁的线程标识  防止误杀
func (r *RedisMutex) setUUID(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key = fmt.Sprintf("uuid_%s", key)
	s := uuid.New().String()
	if err := r.client.Set(key, s, 0).Err(); err != nil {
		return errors.WithStack(err)
	}
	r.m[key] = s
	return nil
}

func (r *RedisMutex) checkUUID(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	key = fmt.Sprintf("uuid_%s", key)

	k, ex := r.m[key]
	if !ex {
		return false
	}
	s := r.client.Get(key).Val()
	if s != k {
		return false
	}
	return true
}

type RedisMutex struct {
	opt    *Option
	client *redis.Client
	m      map[string]string        // 防止误杀
	cl     map[string]chan struct{} // 通知停止心跳
	mu     sync.Mutex
}

func New(fns ...OptionFn) (*RedisMutex, error) {
	opt := Option{}
	for _, v := range fns {
		v(&opt)
	}

	rdOp := &redis.Options{
		Addr:     opt.Uri,
		Password: "",
		DB:       0,
	}
	if opt.Password != nil {
		rdOp.Password = *opt.Password
	}
	if opt.db != nil {
		rdOp.DB = *opt.db
	}

	redisClient := redis.NewClient(rdOp)
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisMutex{
		opt:    &opt,
		client: redisClient,
		m:      map[string]string{},
		cl:     map[string]chan struct{}{},
	}, nil
}

func (r *RedisMutex) Look() (bool, error) {
	return r.HLook(BASE_KEY)
}

func (r *RedisMutex) UnLook() (bool, error) {
	return r.HUnLook(BASE_KEY)
}

func (r *RedisMutex) HLook(key string) (bool, error) {
	if !r.client.SetNX(key, "1", 0).Val() {
		return false, nil
	}

	if r.opt.ExpireTime > 0 {
		if err := r.client.Expire(key, r.opt.ExpireTime).Err(); err != nil {
			return false, errors.WithStack(err)
		}
	}

	cC := make(chan struct{}, 0)

	// set uuid
	if err := r.setUUID(key); err != nil {
		return false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.cl[key] = cC

	go r.heartbeat(key)
	return true, nil
}

func (r *RedisMutex) HUnLook(key string) (bool, error) {
	if !r.checkUUID(key) {
		return false, errors.New("Lock theft")
	}

	if err := r.client.Del(key).Err(); err != nil {
		return false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	c, ex := r.cl[key]
	if !ex {
		return false, errors.New("Lock failure")
	}
	close(c)

	return true, nil
}

func (r *RedisMutex) heartbeat(key string) {
	over := r.cl[key]
	if r.opt.ExpireTime <= 0 {
		return
	}

	// 提前续费
	tim := time.Second * 2
	if r.opt.ExpireTime >= tim {
		tim = time.Second
	} else {
		tim = time.Millisecond * 500
	}
	tm := time.NewTicker(tim)

loop:
	for {
		select {
		case <-tm.C:
			if r.opt.ExpireTime > 0 {
				if err := r.client.Expire(key, r.opt.ExpireTime).Err(); err != nil {
					log.Println(err)
				}
				fmt.Println("续费成功")
			}
		case <-over:
			delete(r.cl, key)
			fmt.Println("关闭成功")
			break loop
		}
	}
}
