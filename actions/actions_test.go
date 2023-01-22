package actions

import (
	"strconv"
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
	"github.com/oxks/mbfl/models"
)

var Incrementator int

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.New("Test_ActionSuite", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

func (as *ActionSuite) CreateUser() *models.User {

	Incrementator++

	uc := strconv.Itoa(Incrementator)
	uc = uc + "_alex@test.com"

	user := &models.User{
		Name:       "Alex",
		Email:      nulls.NewString(uc),
		Provider:   "google",
		ProviderID: "11333_" + uc,
	}

	as.NoError(as.DB.Create(user))

	return user
}

func (as *ActionSuite) Login() *models.User {
	user := as.CreateUser()
	as.Session.Set("current_user_id", user.ID)
	return user
}

func (as *ActionSuite) CreateLink(user *models.User) *models.Link {
	Incrementator++
	link := &models.Link{
		Link:   "http://example.com",
		UserID: user.ID,
	}

	verrs, err := as.DB.ValidateAndCreate(link)
	as.NoError(err)
	as.False(verrs.HasAny())

	return link
}
