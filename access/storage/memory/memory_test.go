package memory_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/seftomsk/abf/access/storage"
	"github.com/seftomsk/abf/access/storage/memory"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}

type StorageSuite struct {
	suite.Suite
	rep       *memory.InMemory
	ctx       context.Context
	ipAddress *storage.IPAddress
}

func (s *StorageSuite) SetupTest() {
	s.ctx = context.Background()
	s.rep = memory.Create()
	s.ipAddress = storage.NewIPAddress("", "", "")
}

func (s *StorageSuite) TestInvalidInitGetErr() {
	rep := memory.InMemory{}

	err := rep.AddToWList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)

	err = rep.AddToBList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)

	err = rep.DeleteFromWhiteList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)

	err = rep.DeleteFromBlackList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)

	_, err = rep.IsInWList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)

	_, err = rep.IsInBList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)
}

func (s *StorageSuite) TestGetDoneFromContextGetErr() {
	s.T().Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := s.rep.AddToWList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.Canceled)

		err = s.rep.AddToBList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.Canceled)

		err = s.rep.DeleteFromWhiteList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.Canceled)

		err = s.rep.DeleteFromBlackList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.Canceled)

		_, err = s.rep.IsInWList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.Canceled)

		_, err = s.rep.IsInBList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.Canceled)
	})
	s.T().Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*0)
		defer cancel()

		err := s.rep.AddToWList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		err = s.rep.AddToBList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		err = s.rep.DeleteFromWhiteList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		err = s.rep.DeleteFromBlackList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		_, err = s.rep.IsInWList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		_, err = s.rep.IsInBList(ctx, s.ipAddress)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)
	})
}

func (s *StorageSuite) TestInvalidEntityGetErr() {
	err := s.rep.AddToWList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidEntity)

	err = s.rep.AddToBList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidEntity)

	err = s.rep.DeleteFromWhiteList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidEntity)

	err = s.rep.DeleteFromBlackList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidEntity)

	_, err = s.rep.IsInWList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidEntity)

	_, err = s.rep.IsInBList(s.ctx, s.ipAddress)
	require.ErrorIs(s.T(), err, storage.ErrInvalidEntity)
}

func (s *StorageSuite) TestAddToWListWithoutErr() {
	err := s.rep.AddToWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))
	require.NoError(s.T(), err)
}

func (s *StorageSuite) TestAddToBListWithoutErr() {
	err := s.rep.AddToBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))
	require.NoError(s.T(), err)
}

func (s *StorageSuite) TestDeleteFromWListWithoutErr() {
	ipAddress := storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128")
	_ = s.rep.AddToWList(s.ctx, ipAddress)
	err := s.rep.DeleteFromWhiteList(s.ctx, ipAddress)
	require.NoError(s.T(), err)
}

func (s *StorageSuite) TestDeleteFromWListNoMaskGetErr() {
	_ = s.rep.AddToWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))

	err := s.rep.DeleteFromWhiteList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.129"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
}

func (s *StorageSuite) TestDeleteFromWListNoIpGetErr() {
	_ = s.rep.AddToWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))

	err := s.rep.DeleteFromWhiteList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.2.0",
		"255.255.255.128"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
}

func (s *StorageSuite) TestDeleteFromBListWithoutErr() {
	ipAddress := storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128")
	_ = s.rep.AddToBList(s.ctx, ipAddress)
	err := s.rep.DeleteFromBlackList(s.ctx, ipAddress)
	require.NoError(s.T(), err)
}

func (s *StorageSuite) TestDeleteFromBListNoMaskGetErr() {
	_ = s.rep.AddToBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))

	err := s.rep.DeleteFromBlackList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.129"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
}

func (s *StorageSuite) TestDeleteFromBListNoIpGetErr() {
	_ = s.rep.AddToBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))

	err := s.rep.DeleteFromBlackList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.2.0",
		"255.255.255.128"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
}

func (s *StorageSuite) TestIsInWListWithoutErr() {
	ipAddress := storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128")
	_ = s.rep.AddToWList(s.ctx, ipAddress)
	exists, err := s.rep.IsInWList(s.ctx, ipAddress)
	require.NoError(s.T(), err)
	require.True(s.T(), exists)
}

func (s *StorageSuite) TestIsInWListNoMaskGetErr() {
	_ = s.rep.AddToWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))
	exists, err := s.rep.IsInWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.129"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
	require.False(s.T(), exists)
}

func (s *StorageSuite) TestIsInWListNoIpGetErr() {
	_ = s.rep.AddToWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))
	exists, err := s.rep.IsInWList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.2.0",
		"255.255.255.128"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
	require.False(s.T(), exists)
}

func (s *StorageSuite) TestIsInBListWithoutErr() {
	ipAddress := storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128")
	_ = s.rep.AddToBList(s.ctx, ipAddress)
	exists, err := s.rep.IsInBList(s.ctx, ipAddress)
	require.NoError(s.T(), err)
	require.True(s.T(), exists)
}

func (s *StorageSuite) TestIsInBListNoMaskGetErr() {
	_ = s.rep.AddToBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))
	exists, err := s.rep.IsInBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.129"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
	require.False(s.T(), exists)
}

func (s *StorageSuite) TestIsInBListNoIpGetErr() {
	_ = s.rep.AddToBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.1.0",
		"255.255.255.128"))
	exists, err := s.rep.IsInBList(s.ctx, storage.NewIPAddress(
		"",
		"192.1.2.0",
		"255.255.255.128"))
	require.ErrorIs(s.T(), err, storage.ErrNotFound)
	require.False(s.T(), exists)
}

//nolint:dupl // Each block code for each case
func (s *StorageSuite) TestMultipleAddToWListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			err := s.rep.AddToWList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			err := s.rep.AddToWList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()
}

//nolint:dupl // Each block code for each case
func (s *StorageSuite) TestMultipleAddToBListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			err := s.rep.AddToBList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			err := s.rep.AddToBList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()
}

//nolint:dupl // Each block code for each case
func (s *StorageSuite) TestMultipleDeleteFromWListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			_ = s.rep.AddToWList(s.ctx, ipAddress)
			err := s.rep.DeleteFromWhiteList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			_ = s.rep.AddToWList(s.ctx, ipAddress)
			err := s.rep.DeleteFromWhiteList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			err := s.rep.DeleteFromWhiteList(s.ctx, ipAddress)
			require.ErrorIs(s.T(), err, storage.ErrNotFound)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			err := s.rep.DeleteFromWhiteList(s.ctx, ipAddress)
			require.ErrorIs(s.T(), err, storage.ErrNotFound)
		}(i)
	}
	wg.Wait()
}

//nolint:dupl // Each block code for each case
func (s *StorageSuite) TestMultipleDeleteFromBListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			_ = s.rep.AddToBList(s.ctx, ipAddress)
			err := s.rep.DeleteFromBlackList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			_ = s.rep.AddToBList(s.ctx, ipAddress)
			err := s.rep.DeleteFromBlackList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			err := s.rep.DeleteFromBlackList(s.ctx, ipAddress)
			require.ErrorIs(s.T(), err, storage.ErrNotFound)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			err := s.rep.DeleteFromBlackList(s.ctx, ipAddress)
			require.ErrorIs(s.T(), err, storage.ErrNotFound)
		}(i)
	}
	wg.Wait()
}

//nolint:dupl // Each block code for each case
func (s *StorageSuite) TestMultipleIsInWListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			_ = s.rep.AddToWList(s.ctx, ipAddress)
			exists, err := s.rep.IsInWList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			_ = s.rep.AddToWList(s.ctx, ipAddress)
			exists, err := s.rep.IsInWList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			exists, err := s.rep.IsInWList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			exists, err := s.rep.IsInWList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	wg.Wait()
}

//nolint:dupl // Each block code for each case
func (s *StorageSuite) TestMultipleIsInBListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			_ = s.rep.AddToBList(s.ctx, ipAddress)
			exists, err := s.rep.IsInBList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			_ = s.rep.AddToBList(s.ctx, ipAddress)
			exists, err := s.rep.IsInBList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.128")
			exists, err := s.rep.IsInBList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			ipAddress := storage.NewIPAddress(
				"",
				fmt.Sprintf("192.1.%v.0", i),
				"255.255.255.129")
			exists, err := s.rep.IsInBList(s.ctx, ipAddress)
			require.NoError(s.T(), err)
			require.True(s.T(), exists)
		}(i)
	}
	wg.Wait()
}
