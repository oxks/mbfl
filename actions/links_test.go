package actions

import (
	"github.com/oxks/mbfl/models"
)

func (as *ActionSuite) Test_LinksResource_List() {

	user := as.Login()
	link := as.CreateLink(user)

	res := as.HTML("/links/").Get()
	as.Equal(200, res.Code)

	as.Contains(res.Body.String(), link.Link)

}

func (as *ActionSuite) Test_LinksResource_List_Not_Owner() {
	u1 := as.CreateUser()
	link := as.CreateLink(u1)

	as.Login()

	res := as.HTML("/links").Get()
	as.Equal(200, res.Code)

	as.NotContains(res.Body.String(), link.Link)

}

func (as *ActionSuite) Test_LinksResource_Show() {
	user := as.Login()
	link := as.CreateLink(user)

	res := as.HTML("/links/%s", link.ID).Get()
	as.Equal(200, res.Code)

	as.Contains(res.Body.String(), link.Link)

}

func (as *ActionSuite) Test_LinksResource_Show_Not_Owner() {
	u1 := as.CreateUser()
	link := as.CreateLink(u1)

	as.Login()

	res := as.HTML("/links/%s", link.ID).Get()
	as.Equal(404, res.Code)

}

func (as *ActionSuite) Test_LinksResource_New() {
	as.Login()
	res := as.HTML("/links/new").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "New Link")
}

func (as *ActionSuite) Test_LinksResource_Create() {
	user := as.Login()
	link := &models.Link{
		Link: "http://gobuffalo.io",
	}
	res := as.HTML("/links").Post(link)
	as.Equal(302, res.Code)

	l := &models.Link{}
	as.NoError(as.DB.First(l))
	as.Equal(link.Link, l.Link)
	as.Equal(user.ID, l.UserID)
	as.NotZero(l.Code)

}

func (as *ActionSuite) Test_LinksResource_Edit() {
	user := as.Login()
	link := as.CreateLink(user)

	res := as.HTML("/links/%s/edit", link.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Edit Link")

}

func (as *ActionSuite) Test_LinksResource_Edit_Not_Owner() {
	u1 := as.CreateUser()
	link := as.CreateLink(u1)

	as.Login()

	res := as.HTML("/links/%s/edit", link.ID).Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_LinksResource_Update() {

	user := as.Login()
	link := as.CreateLink(user)
	pl := link.Link
	link.Link = "http://something.new"
	res := as.HTML("/links/%s", link.ID).Put(link)
	as.Equal(302, res.Code)

	l := &models.Link{}
	as.NoError(as.DB.First(l))
	as.NotEqual(pl, l.Link)
	as.Equal(user.ID, l.UserID)
	as.NotZero(l.Code)
}

func (as *ActionSuite) Test_LinksResource_Update_Not_Owner() {

	u1 := as.CreateUser()
	link := as.CreateLink(u1)

	as.Login()

	link.Link = "http://something.new"

	res := as.HTML("/links/%s", link.ID).Put(link)
	as.Equal(404, res.Code)

}

func (as *ActionSuite) Test_LinksResource_Destroy() {
	user := as.Login()
	link := as.CreateLink(user)
	res := as.HTML("/links/%s", link.ID).Delete()
	as.Equal(302, res.Code)
	as.Equal("/links/", res.Location())

}

func (as *ActionSuite) Test_LinksResource_Destroy_Not_Owner() {
	u1 := as.CreateUser()
	link := as.CreateLink(u1)

	as.Login()

	res := as.HTML("/links/%s", link.ID).Delete()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_Redirector() {
	link := as.CreateLink(as.CreateUser())
	// count, err := as.DB.Count("clicks")
	// as.NoError(err)
	// as.Equal(0, count)

	res := as.HTML(link.Code).Get()
	// as.Equal(302, res.Code)
	as.Equal("/"+link.Code+"/", res.Location())

	// as.Equal(link.Link, res.Location())

	// count, err = as.DB.Count("clicks")
	// as.NoError(err)
	// as.Equal(1, count)
}
