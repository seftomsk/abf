package access_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"

	"github.com/seftomsk/abf/access"

	"github.com/seftomsk/abf/access/storage"

	"github.com/seftomsk/abf/access/storage/memory"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(AccessSuite))
}

type AccessSuite struct {
	suite.Suite
	ctx      context.Context
	access   *access.IPAccess
	emptyDTO access.IPDTO
	validDTO access.IPDTO
}

func (s *AccessSuite) SetupTest() {
	s.ctx = context.Background()
	s.access = access.NewIPAccess(memory.Create())
	s.emptyDTO = access.IPDTO{}
	s.validDTO = access.IPDTO{IP: "192.1.1.0/25"}
}

func (s *AccessSuite) TestInvalidStorageGetErr() {
	a := access.IPAccess{}

	err := a.AddToWList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrInvalidStorage)

	err = a.AddToBList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrInvalidStorage)

	err = a.DeleteFromWList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrInvalidStorage)

	err = a.DeleteFromBList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrInvalidStorage)

	ok, err := a.IsInWList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrInvalidStorage)
	require.False(s.T(), ok)

	ok, err = a.IsInBList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrInvalidStorage)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestGetDoneFromContextGetErr() {
	s.T().Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := s.access.AddToWList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.Canceled)

		err = s.access.AddToBList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.Canceled)

		err = s.access.DeleteFromWList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.Canceled)

		err = s.access.DeleteFromBList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.Canceled)

		_, err = s.access.IsInWList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.Canceled)

		_, err = s.access.IsInBList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.Canceled)
	})
	s.T().Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*0)
		defer cancel()

		err := s.access.AddToWList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		err = s.access.AddToBList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		err = s.access.DeleteFromWList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		err = s.access.DeleteFromBList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		_, err = s.access.IsInWList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)

		_, err = s.access.IsInBList(ctx, s.emptyDTO)
		require.ErrorIs(s.T(), err, context.DeadlineExceeded)
	})
}

func (s *AccessSuite) TestEmptyIpGetErr() {
	err := s.access.AddToWList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrEmptyIP)

	err = s.access.AddToBList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrEmptyIP)

	err = s.access.DeleteFromWList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrEmptyIP)

	err = s.access.DeleteFromBList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrEmptyIP)

	ok, err := s.access.IsInWList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrEmptyIP)
	require.False(s.T(), ok)

	ok, err = s.access.IsInBList(s.ctx, s.emptyDTO)
	require.ErrorIs(s.T(), err, access.ErrEmptyIP)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestInvalidIpGetErr() {
	dto := access.IPDTO{IP: "a"}
	var e *access.ErrParseIP

	err := s.access.AddToWList(s.ctx, dto)
	require.ErrorAs(s.T(), err, &e)

	err = s.access.AddToBList(s.ctx, dto)
	require.ErrorAs(s.T(), err, &e)

	err = s.access.DeleteFromWList(s.ctx, dto)
	require.ErrorAs(s.T(), err, &e)

	err = s.access.DeleteFromBList(s.ctx, dto)
	require.ErrorAs(s.T(), err, &e)

	ok, err := s.access.IsInWList(s.ctx, dto)
	require.ErrorAs(s.T(), err, &e)
	require.False(s.T(), ok)

	ok, err = s.access.IsInBList(s.ctx, dto)
	require.ErrorAs(s.T(), err, &e)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestAddToWListWithoutErr() {
	err := s.access.AddToWList(s.ctx, s.validDTO)
	require.NoError(s.T(), err)
}

func (s *AccessSuite) TestAddToWListGetErr() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	m := NewMockIStorage(ctrl)
	m.EXPECT().
		AddToWList(gomock.Any(), gomock.Any()).
		Return(storage.ErrInvalidInitialization)
	a := access.NewIPAccess(m)
	err := a.AddToWList(s.ctx, s.validDTO)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)
}

func (s *AccessSuite) TestAddToBListWithoutErr() {
	err := s.access.AddToBList(s.ctx, s.validDTO)
	require.NoError(s.T(), err)
}

func (s *AccessSuite) TestAddToBListGetErr() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	m := NewMockIStorage(ctrl)
	m.EXPECT().
		AddToBList(gomock.Any(), gomock.Any()).
		Return(storage.ErrInvalidInitialization)
	a := access.NewIPAccess(m)
	err := a.AddToBList(s.ctx, s.validDTO)
	require.ErrorIs(s.T(), err, storage.ErrInvalidInitialization)
}

func (s *AccessSuite) TestDeleteFromWListWithoutErr() {
	_ = s.access.AddToWList(s.ctx, s.validDTO)
	err := s.access.DeleteFromWList(s.ctx, s.validDTO)
	require.NoError(s.T(), err)
}

func (s *AccessSuite) TestDeleteFromWListNotMaskGetErr() {
	err := s.access.DeleteFromWList(s.ctx, s.validDTO)
	require.ErrorIs(s.T(), err, access.ErrNotFound)
}

func (s *AccessSuite) TestDeleteFromWListNotIpGetErr() {
	_ = s.access.AddToWList(s.ctx, s.validDTO)
	err := s.access.DeleteFromWList(s.ctx, access.IPDTO{IP: "192.1.2.0/25"})
	require.ErrorIs(s.T(), err, access.ErrNotFound)
}

func (s *AccessSuite) TestDeleteFromWListGetErr() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	m := NewMockIStorage(ctrl)
	m.EXPECT().
		DeleteFromWList(gomock.Any(), gomock.Any()).
		Return(storage.ErrInvalidEntity)
	a := access.NewIPAccess(m)
	err := a.DeleteFromWList(s.ctx, s.validDTO)
	require.Error(s.T(), err)
}

func (s *AccessSuite) TestDeleteFromBListWithoutErr() {
	_ = s.access.AddToBList(s.ctx, s.validDTO)
	err := s.access.DeleteFromBList(s.ctx, s.validDTO)
	require.NoError(s.T(), err)
}

func (s *AccessSuite) TestDeleteFromBListNotMaskGetErr() {
	err := s.access.DeleteFromBList(s.ctx, s.validDTO)
	require.ErrorIs(s.T(), err, access.ErrNotFound)
}

func (s *AccessSuite) TestDeleteFromBListNotIpGetErr() {
	_ = s.access.AddToBList(s.ctx, s.validDTO)
	err := s.access.DeleteFromBList(s.ctx, access.IPDTO{IP: "192.1.2.0/25"})
	require.ErrorIs(s.T(), err, access.ErrNotFound)
}

func (s *AccessSuite) TestDeleteFromBListGetErr() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	m := NewMockIStorage(ctrl)
	m.EXPECT().
		DeleteFromBList(gomock.Any(), gomock.Any()).
		Return(storage.ErrInvalidEntity)
	a := access.NewIPAccess(m)
	err := a.DeleteFromBList(s.ctx, s.validDTO)
	require.Error(s.T(), err)
}

func (s *AccessSuite) TestIsInWListWithoutErr() {
	_ = s.access.AddToWList(s.ctx, s.validDTO)
	ok, err := s.access.IsInWList(s.ctx, s.validDTO)
	require.NoError(s.T(), err)
	require.True(s.T(), ok)
}

func (s *AccessSuite) TestIsInWListNoMaskGetErr() {
	ok, err := s.access.IsInWList(s.ctx, s.validDTO)
	require.ErrorIs(s.T(), err, access.ErrNotFound)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestIsInWListNoIpGetErr() {
	_ = s.access.AddToWList(s.ctx, s.validDTO)
	ok, err := s.access.IsInWList(s.ctx, access.IPDTO{IP: "192.1.2.0/25"})
	require.ErrorIs(s.T(), err, access.ErrNotFound)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestIsInWListGetErr() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	m := NewMockIStorage(ctrl)
	m.EXPECT().
		IsInWList(gomock.Any(), gomock.Any()).
		Return(false, storage.ErrInvalidEntity)
	a := access.NewIPAccess(m)
	_, err := a.IsInWList(s.ctx, s.validDTO)
	require.Error(s.T(), err)
}

func (s *AccessSuite) TestIsInBListWithoutErr() {
	_ = s.access.AddToBList(s.ctx, s.validDTO)
	ok, err := s.access.IsInBList(s.ctx, s.validDTO)
	require.NoError(s.T(), err)
	require.True(s.T(), ok)
}

func (s *AccessSuite) TestIsInBListNoMaskGetErr() {
	ok, err := s.access.IsInBList(s.ctx, s.validDTO)
	require.ErrorIs(s.T(), err, access.ErrNotFound)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestIsInBListNoIpGetErr() {
	_ = s.access.AddToBList(s.ctx, s.validDTO)
	ok, err := s.access.IsInBList(s.ctx, access.IPDTO{IP: "192.1.2.0/25"})
	require.ErrorIs(s.T(), err, access.ErrNotFound)
	require.False(s.T(), ok)
}

func (s *AccessSuite) TestIsInBListGetErr() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
	m := NewMockIStorage(ctrl)
	m.EXPECT().
		IsInBList(gomock.Any(), gomock.Any()).
		Return(false, storage.ErrInvalidEntity)
	a := access.NewIPAccess(m)
	_, err := a.IsInBList(s.ctx, s.validDTO)
	require.Error(s.T(), err)
}

func (s *AccessSuite) TestMultipleAddToWListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			err := s.access.AddToWList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			err := s.access.AddToWList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()
}

func (s *AccessSuite) TestMultipleAddToBListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			err := s.access.AddToBList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			err := s.access.AddToBList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()
}

func (s *AccessSuite) TestMultipleDeleteFromWListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			_ = s.access.AddToWList(s.ctx, dto)
			err := s.access.DeleteFromWList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			_ = s.access.AddToWList(s.ctx, dto)
			err := s.access.DeleteFromWList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			err := s.access.DeleteFromWList(s.ctx, dto)
			require.ErrorIs(s.T(), err, access.ErrNotFound)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			err := s.access.DeleteFromWList(s.ctx, dto)
			require.ErrorIs(s.T(), err, access.ErrNotFound)
		}(i)
	}
	wg.Wait()
}

func (s *AccessSuite) TestMultipleDeleteFromBListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			_ = s.access.AddToBList(s.ctx, dto)
			err := s.access.DeleteFromBList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			_ = s.access.AddToBList(s.ctx, dto)
			err := s.access.DeleteFromBList(s.ctx, dto)
			require.NoError(s.T(), err)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			err := s.access.DeleteFromBList(s.ctx, dto)
			require.ErrorIs(s.T(), err, access.ErrNotFound)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			err := s.access.DeleteFromBList(s.ctx, dto)
			require.ErrorIs(s.T(), err, access.ErrNotFound)
		}(i)
	}
	wg.Wait()
}

func (s *AccessSuite) TestMultipleIsInWListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			_ = s.access.AddToWList(s.ctx, dto)
			ok, err := s.access.IsInWList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			_ = s.access.AddToWList(s.ctx, dto)
			ok, err := s.access.IsInWList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			ok, err := s.access.IsInWList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			ok, err := s.access.IsInWList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	wg.Wait()
}

func (s *AccessSuite) TestMultipleIsInBListWithoutErr() {
	wg := &sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			_ = s.access.AddToBList(s.ctx, dto)
			ok, err := s.access.IsInBList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			_ = s.access.AddToBList(s.ctx, dto)
			ok, err := s.access.IsInBList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/24", i)}
			ok, err := s.access.IsInBList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			dto := access.IPDTO{IP: fmt.Sprintf("192.1.%v.0/25", i)}
			ok, err := s.access.IsInBList(s.ctx, dto)
			require.NoError(s.T(), err)
			require.True(s.T(), ok)
		}(i)
	}
	wg.Wait()
}
