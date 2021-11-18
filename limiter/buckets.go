package limiter

type Buckets struct {
	collection []IBucket
}

func (bs *Buckets) AddTokens() {
	for _, bucket := range bs.collection {
		bucket.AddTokens()
	}
}

func (bs *Buckets) DeleteToken() {
	for _, bucket := range bs.collection {
		bucket.DeleteToken()
	}
}

func (bs *Buckets) CountAvailableTokens() int {
	var tokens int
	for _, bucket := range bs.collection {
		if bucket.CountAvailableTokens() > 0 {
			tokens += bucket.CountAvailableTokens()
		}
	}

	return tokens
}

func (bs *Buckets) CheckTokensExist() bool {
	for _, bucket := range bs.collection {
		if bucket.CountAvailableTokens() <= 0 {
			return false
		}
	}

	return true
}

func (bs *Buckets) ClearBucket() {
	for _, bucket := range bs.collection {
		bucket.ClearBucket()
	}
}
