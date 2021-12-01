package limiter

type MultiLimiter struct {
	login    *Limiter
	password *Limiter
	ip       *Limiter
}

func NewMultiLimiter(login, password, ip *Limiter) *MultiLimiter {
	return &MultiLimiter{
		login:    login,
		password: password,
		ip:       ip,
	}
}

func (ml *MultiLimiter) GetBucket(login, password, ip string) IBucket {
	lBucket := ml.GetLoginBucket(login)
	pBucket := ml.GetPasswordBucket(password)
	iBucket := ml.GetIPBucket(ip)

	return &Buckets{collection: []IBucket{lBucket, pBucket, iBucket}}
}

func (ml *MultiLimiter) GetLoginBucket(key string) IBucket {
	return ml.login.GetBucket(key)
}

func (ml *MultiLimiter) GetPasswordBucket(key string) IBucket {
	return ml.password.GetBucket(key)
}

func (ml *MultiLimiter) GetIPBucket(key string) IBucket {
	return ml.ip.GetBucket(key)
}
