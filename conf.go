package redis_mutex

import (
	"fmt"
	"time"
)

const BASE_KEY = `BASE_KEY`

type Option struct {
	Uri        string
	ExpireTime time.Duration

	Password *string
	db       *int
}

type OptionFn func(o *Option)

func SetUri(host string, port int) OptionFn {
	return func(o *Option) {
		o.Uri = fmt.Sprintf("%s:%d", host, port)
		o.ExpireTime = time.Second * 5
	}
}

func SetPassword(password string) OptionFn {
	return func(o *Option) {
		o.Password = &password
	}
}

func SetDB(db int) OptionFn {
	return func(o *Option) {
		o.db = &db
	}
}

func SetExpireTime(time time.Duration) OptionFn {
	return func(o *Option) {
		o.ExpireTime = time
	}
}
