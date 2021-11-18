package access_test

import (
	"fmt"
	"github.com/seftomsk/abf/access"
	"github.com/seftomsk/abf/access/storage"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AccessSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AccessSuite))
}

func (a *AccessSuite) TestAddToWhiteList() {
	//ipAccess := access.IPAccess{}
	//require.NoError(a.T(), ipAccess.AddToWhiteList("192.1.1.0/25"))
}

func (a *AccessSuite) TestGetAll() {
	memory := storage.NewInMemory()
	ipAccess := access.NewIPAccess(memory)
	err := ipAccess.AddToWhiteList("192.1.1.0/25")
	require.NoError(a.T(), err)

	err = ipAccess.AddToWhiteList("192.1.1.1/25")
	require.NoError(a.T(), err)

	err = ipAccess.AddToWhiteList("192.1.1.2/25")
	require.NoError(a.T(), err)

	err = ipAccess.AddToBlackList("192.1.1.1/25")
	require.NoError(a.T(), err)

	err = ipAccess.AddToBlackList("192.1.1.2/25")
	require.NoError(a.T(), err)

	fmt.Println(ipAccess.GetAll())
}