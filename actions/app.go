// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {

		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_saloon_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// middleware to set a current_user_id session
		app.Use(SetCurrentUser)
		app.Use(SetCurrentForum)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)

		app.ServeFiles("/assets", assetsBox)

		auth := app.Group("/users")
		auth.GET("/register", UsersRegisterGet)
		auth.POST("/register", UsersRegisterPost)
		auth.GET("/login", UsersLoginGet)
		auth.POST("/login", UsersLoginPost)
		auth.GET("/logout", UsersLogout)
		auth.GET("/settings", UserRequired(UsersSettings))
		auth.GET("/show", UserRequired(UsersShow))
		auth.GET("/settings/add-subscription/{cid}", UserRequired(UsersSettingsAddSubscription))
		auth.GET("/settings/rm-subscription/{cid}", UserRequired(UsersSettingsRemoveSubscription))
		auth.POST("/settings/update-avatar", UserRequired(UsersSettingsUpdateAvatar))
		auth.POST("/settings/update-name", UserRequired(UsersSettingsUpdateName))
		auth.POST("/settings/update-email", UserRequired(UsersSettingsUpdateEmail))
		auth.POST("/settings/update-password", UserRequired(UsersSettingsUpdatePassword))

		catGroup := app.Group("/categories")
		catGroup.Use(UserRequired)
		catGroup.GET("/index", CategoriesIndex)
		catGroup.GET("/create", AdminRequired(CategoriesCreateGet))
		catGroup.POST("/create", AdminRequired(CategoriesCreatePost))
		catGroup.GET("/detail/{cid}", CategoriesDetail)

		topicGroup := app.Group("/topics")
		topicGroup.Use(UserRequired)
		topicGroup.GET("/index", TopicsIndex)
		topicGroup.GET("/detail/{tid}", TopicsDetail)
		topicGroup.GET("/create", TopicsCreateGet)
		topicGroup.POST("/create", TopicsCreatePost)
		topicGroup.GET("/delete", TopicsDelete)
		topicGroup.GET("/edit", TopicsEditGet)
		topicGroup.POST("/edit", TopicsEditPost)
		topicGroup.GET("/add-subscriber/{tid}", UserRequired(TopicsAddSubscriber))
		topicGroup.GET("/rm-subscriber/{tid}", UserRequired(TopicsRemoveSubscriber))

		replyGroup := app.Group("/replies")
		replyGroup.Use(UserRequired)
		replyGroup.GET("/create", RepliesCreateGet)
		replyGroup.POST("/create", RepliesCreatePost)
		replyGroup.GET("/edit", RepliesEdit)
		replyGroup.GET("/delete", RepliesDelete)
		replyGroup.GET("/detail", RepliesDetail)

		app.GET("/search", UserRequired(Search))
	}

	return app
}
