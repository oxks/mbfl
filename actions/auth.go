package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/markbates/going/defaults"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/oxks/mbfl/models"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	hosterAppHost := App().Host
	if os.Getenv("HOSTER_APP_HOST") != nil {
		hosterAppHost = os.Getenv("HOSTER_APP_HOST")
	}
	goth.UseProviders(

		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", hosterAppHost, "/auth/google/callback")),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", hosterAppHost, "/auth/github/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	gu, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("provider = ? and provider_id = ?", gu.Provider, gu.UserID)
	exists, err := q.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}
	u := &models.User{}
	if exists {
		if err = q.First(u); err != nil {
			return errors.WithStack(err)
		}
	}
	u.Name = defaults.String(gu.Name, gu.NickName)
	u.Provider = gu.Provider
	u.ProviderID = gu.UserID
	u.Email = nulls.NewString(gu.Email)
	if err = tx.Save(u); err != nil {
		return errors.WithStack(err)
	}

	c.Session().Set("current_user_id", u.ID)
	if err = c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "You have been logged in")
	return c.Redirect(302, "/links")
}

func AuthDestroy(c buffalo.Context) error {
	c.Flash().Add("danger", "my destr")
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out")
	return c.Redirect(302, "/")
}

func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			if err := tx.Find(u, uid); err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
