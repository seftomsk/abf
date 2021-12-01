package limiter_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/seftomsk/abf/limiter"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(LimiterSuite))
	suite.Run(t, new(LimitersSuite))
}

type LimiterSuite struct {
	suite.Suite
	limiter   *limiter.Limiter
	bucket    limiter.IBucket
	firstKey  string
	secondKey string
}

func (suite *LimiterSuite) SetupTest() {
	suite.limiter = limiter.NewLimiter(4, time.Second)
	suite.firstKey = "user"
	suite.secondKey = "admin"
	suite.bucket = suite.limiter.GetBucket(suite.firstKey)
}

type LimitersSuite struct {
	suite.Suite
	limiter *limiter.MultiLimiter
	bucket  limiter.IBucket
}

func (suite *LimitersSuite) SetupTest() {
	loginLimiter := limiter.NewLimiter(2, time.Second)
	passwordLimiter := limiter.NewLimiter(2, time.Second)
	ipLimiter := limiter.NewLimiter(2, time.Second)

	multiLimiter := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)

	bucket := multiLimiter.GetBucket(
		"user",
		"0",
		"127.0.0.1")

	suite.limiter = multiLimiter
	suite.bucket = bucket
}

func (suite *LimiterSuite) TestTokensExist() {
	require.True(suite.T(), suite.bucket.CheckTokensExist())
}

func (suite *LimiterSuite) TestTokensAreNotExist() {
	loginLimiter := limiter.NewLimiter(0, time.Microsecond)
	bucket := loginLimiter.GetBucket(suite.firstKey)
	require.False(suite.T(), bucket.CheckTokensExist())
}

func (suite *LimiterSuite) TestAvailableTokens() {
	require.Equal(suite.T(), 4, suite.bucket.CountAvailableTokens())
}

func (suite *LimiterSuite) TestDeleteToken() {
	suite.bucket.DeleteToken()
	suite.bucket.DeleteToken()
	require.Equal(suite.T(), 2, suite.bucket.CountAvailableTokens())
}

func (suite *LimiterSuite) TestAddTokens() {
	loginLimiter := limiter.NewLimiter(4, time.Microsecond)
	bucket := loginLimiter.GetBucket(suite.firstKey)
	bucket.DeleteToken()
	bucket.DeleteToken()
	time.Sleep(time.Millisecond)
	bucket.AddTokens()
	require.Equal(suite.T(), 4, bucket.CountAvailableTokens())
}

func (suite *LimiterSuite) TestClearBucket() {
	suite.bucket.ClearBucket()
	require.Equal(suite.T(), 0, suite.bucket.CountAvailableTokens())
}

func (suite *LimiterSuite) TestAddTokensBeforeDuration() {
	loginLimiter := limiter.NewLimiter(4, time.Minute)
	bucket := loginLimiter.GetBucket(suite.firstKey)
	bucket.DeleteToken()
	bucket.DeleteToken()
	bucket.AddTokens()
	require.Equal(suite.T(), 2, bucket.CountAvailableTokens())
}

func (suite *LimiterSuite) TestMultipleAccessingToAvailableTokens() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			bucket := suite.limiter.GetBucket(suite.firstKey)
			require.Equal(suite.T(), 4, bucket.CountAvailableTokens())
		}()
		go func() {
			defer wg.Done()
			bucket := suite.limiter.GetBucket(suite.secondKey)
			require.Equal(suite.T(), 4, bucket.CountAvailableTokens())
		}()
	}
	wg.Wait()
}

func (suite *LimiterSuite) TestMultipleDeletingToken() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			bucket := suite.limiter.GetBucket(suite.firstKey)
			bucket.DeleteToken()
		}()

		go func() {
			defer wg.Done()
			bucket := suite.limiter.GetBucket(suite.secondKey)
			bucket.DeleteToken()
		}()
	}
	wg.Wait()
	availableTokens := suite.limiter.
		GetBucket(suite.firstKey).
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)

	availableTokens = suite.limiter.
		GetBucket(suite.secondKey).
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)
}

func (suite *LimiterSuite) TestMultipleAddingTokens() {
	loginLimiter := limiter.NewLimiter(4, time.Microsecond)
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			bucket := loginLimiter.GetBucket(suite.firstKey)
			bucket.DeleteToken()
			bucket.DeleteToken()
			time.Sleep(time.Millisecond)
			bucket.AddTokens()
		}()
		go func() {
			defer wg.Done()
			bucket := loginLimiter.GetBucket(suite.secondKey)
			bucket.DeleteToken()
			bucket.DeleteToken()
			time.Sleep(time.Millisecond)
			bucket.AddTokens()
		}()
	}
	wg.Wait()
	availableTokens := loginLimiter.
		GetBucket(suite.firstKey).
		CountAvailableTokens()
	require.Equal(suite.T(), 4, availableTokens)

	availableTokens = loginLimiter.
		GetBucket(suite.secondKey).
		CountAvailableTokens()
	require.Equal(suite.T(), 4, availableTokens)
}

// Limiters.
func (suite *LimitersSuite) TestTokensExist() {
	require.True(suite.T(), suite.bucket.CheckTokensExist())
}

func (suite *LimitersSuite) TestTokensAreNotExistInAllBuckets() {
	loginLimiter := limiter.NewLimiter(0, time.Second*10)
	passwordLimiter := limiter.NewLimiter(0, time.Second*12)
	ipLimiter := limiter.NewLimiter(0, time.Second*14)

	multiLimiter := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)

	buckets := multiLimiter.GetBucket(
		"user",
		"0",
		"127.0.0.1")

	require.False(suite.T(), buckets.CheckTokensExist())
}

func (suite *LimitersSuite) TestTokensAreNotExistInAnyBucket() {
	loginLimiter := limiter.NewLimiter(10, time.Second*10)
	passwordLimiter := limiter.NewLimiter(0, time.Second*12)
	ipLimiter := limiter.NewLimiter(20, time.Second*14)

	multiLimiter := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)

	buckets := multiLimiter.GetBucket(
		"user",
		"0",
		"127.0.0.1")

	require.False(suite.T(), buckets.CheckTokensExist())
}

func (suite *LimitersSuite) TestDeleteToken() {
	suite.bucket.DeleteToken() // Delete one token from each bucket (Total -3)
	require.Equal(suite.T(), 3, suite.bucket.CountAvailableTokens())
}

func (suite *LimitersSuite) TestAddTokens() {
	loginLimiter := limiter.NewLimiter(1, time.Microsecond)
	passwordLimiter := limiter.NewLimiter(1, time.Microsecond)
	ipLimiter := limiter.NewLimiter(1, time.Microsecond)

	multiLimiter := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)

	buckets := multiLimiter.GetBucket(
		"user",
		"0",
		"127.0.0.1")

	buckets.DeleteToken()
	time.Sleep(time.Millisecond)
	buckets.AddTokens()
	require.Equal(suite.T(), 3, buckets.CountAvailableTokens())
}

func (suite *LimitersSuite) TestClearBucket() {
	suite.bucket.ClearBucket()
	require.Equal(suite.T(), 0, suite.bucket.CountAvailableTokens())
}

func (suite *LimitersSuite) TestAddTokensBeforeDuration() {
	loginLimiter := limiter.NewLimiter(2, time.Minute)
	passwordLimiter := limiter.NewLimiter(2, time.Minute)
	ipLimiter := limiter.NewLimiter(2, time.Minute)

	multiLimiter := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)

	buckets := multiLimiter.GetBucket(
		"user",
		"0",
		"127.0.0.1")

	buckets.DeleteToken()
	buckets.DeleteToken()
	buckets.AddTokens()
	require.Equal(suite.T(), 0, buckets.CountAvailableTokens())
}

func (suite *LimitersSuite) TestMultipleAccessingToAvailableTokens() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			buckets := suite.limiter.GetBucket(
				"user",
				"0",
				"127.0.0.1")
			require.Equal(suite.T(), 6, buckets.CountAvailableTokens())
		}()
		go func() {
			defer wg.Done()
			buckets := suite.limiter.GetBucket(
				"user2",
				"1",
				"127.0.0.2")
			require.Equal(suite.T(), 6, buckets.CountAvailableTokens())
		}()
	}
	wg.Wait()
}

func (suite *LimitersSuite) TestMultipleDeletingToken() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			buckets := suite.limiter.GetBucket(
				"user",
				"0",
				"127.0.0.1")
			buckets.DeleteToken()
		}()

		go func() {
			defer wg.Done()
			buckets := suite.limiter.GetBucket(
				"user2",
				"1",
				"127.0.0.2")
			buckets.DeleteToken()
		}()
	}
	wg.Wait()

	availableTokens := suite.limiter.
		GetLoginBucket("user").
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)

	availableTokens = suite.limiter.
		GetLoginBucket("user2").
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)

	availableTokens = suite.limiter.
		GetPasswordBucket("0").
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)

	availableTokens = suite.limiter.
		GetPasswordBucket("1").
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)

	availableTokens = suite.limiter.
		GetIPBucket("127.0.0.1").
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)

	availableTokens = suite.limiter.
		GetIPBucket("127.0.0.2").
		CountAvailableTokens()
	require.Equal(suite.T(), 0, availableTokens)
}

func (suite *LimitersSuite) TestMultipleAddingTokens() {
	loginLimiter := limiter.NewLimiter(2, time.Microsecond)
	passwordLimiter := limiter.NewLimiter(2, time.Microsecond)
	ipLimiter := limiter.NewLimiter(2, time.Microsecond)

	multiLimiter := limiter.NewMultiLimiter(
		loginLimiter,
		passwordLimiter,
		ipLimiter)

	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			buckets := multiLimiter.GetBucket(
				"user",
				"0",
				"127.0.0.1")
			buckets.DeleteToken()
			buckets.DeleteToken()
			time.Sleep(time.Millisecond)
			buckets.AddTokens()
		}()
		go func() {
			defer wg.Done()
			buckets := multiLimiter.GetBucket(
				"user2",
				"1",
				"127.0.0.2")
			buckets.DeleteToken()
			buckets.DeleteToken()
			time.Sleep(time.Millisecond)
			buckets.AddTokens()
		}()
	}
	wg.Wait()

	availableTokens := multiLimiter.
		GetLoginBucket("user").
		CountAvailableTokens()
	require.Equal(suite.T(), 2, availableTokens)

	availableTokens = multiLimiter.
		GetLoginBucket("user2").
		CountAvailableTokens()
	require.Equal(suite.T(), 2, availableTokens)

	availableTokens = multiLimiter.
		GetPasswordBucket("0").
		CountAvailableTokens()
	require.Equal(suite.T(), 2, availableTokens)

	availableTokens = multiLimiter.
		GetPasswordBucket("1").
		CountAvailableTokens()
	require.Equal(suite.T(), 2, availableTokens)

	availableTokens = multiLimiter.
		GetIPBucket("127.0.0.1").
		CountAvailableTokens()
	require.Equal(suite.T(), 2, availableTokens)

	availableTokens = multiLimiter.
		GetIPBucket("127.0.0.2").
		CountAvailableTokens()
	require.Equal(suite.T(), 2, availableTokens)
}
