package app

import (
	"goshop/app/database"
	"labix.org/v2/mgo"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/gzip"
	"github.com/codegangsta/martini-contrib/render"
	"reflect"
)

type Shop struct {
	databaseSession *mgo.Session
	Server *martini.ClassicMartini
	Categories []Category
	Routes []interface{}
}

type Page struct {
	Title, Description string
	Path, Template, Slug string
	HttpCode int
	HttpMethod string
}

type InternalRoute struct {
	Path string
	HttpMethod string
	Run func(params martini.Params, r render.Render)
}

type Category struct {
	Slug, Title, Description string
	Products []Product
}

type Product struct {
	Slug, Title, Description, Image string
	Price float64
	Off float64
	Category Category
}

// Start server and connection to database
func (app *Shop) Start() *Shop {
	// don't connect if is connected
	if app.databaseSession == nil {
		app.databaseSession = database.Connect("localhost", "heramodas")
	}

	// create martini instance
	app.Server = martini.Classic()

	// gzip all requests
	app.Server.Use(gzip.All())

	// static files
	app.Server.Use(martini.Static("assets"))

	return app
}

// Get category by slug in database
func (app *Shop) GetCategory(slug string, products bool) (Category, error) {
	var c Category

	// find category
	err := database.Find("categories", database.M{"slug": slug}).One(&c)

	if products == true && err == nil {
		// look for category products
		database.Find("products", database.M{"category": slug}).All(&c.Products)
	}

	return c, err
}

// Get product details in database
func (app *Shop) GetProduct(category string, slug string) (Product, error) {
	var p Product

	// finc product
	err := database.Find("products", database.M{"category": category, "slug": slug}).One(&p)

	return p, err
}

// Define routes
func (app *Shop) Route() *Shop {
	// Register routes
	for i := 0; i < len(app.Routes); i++ {
		page := app.Routes[i]

		fields := reflect.Indirect(reflect.ValueOf(page))

		switch fields.FieldByName("HttpMethod").String() {
		case "GET":
			switch page.(type) {
				case Page:
					// convert route to page struct
					route := page.(Page)

					app.Server.Get(route.Path, func(r render.Render) {
						r.HTML(route.HttpCode, route.Template, route)
					})
				case InternalRoute:
					// convert route to InternalRoute struct
					r := page.(InternalRoute)

					app.Server.Get(r.Path, r.Run)
			}
		}
	}

	return app
}

// Category router
func (app *Shop) RouteCategory(params martini.Params, r render.Render) {
	slug := params["category"]

	if slug == "" {
		// if there is no parameter, render 404 error
		renderError(404, r)
	} else {
		// look for the category and their products
		if category, err := app.GetCategory(slug, true); err == nil {
			// render category page
			r.HTML(200, "category", category)
		} else {
			// category not found
			renderError(404, r)
		}
	}
}

// Product details router
func (app *Shop) RouteProduct(params martini.Params, r render.Render) {
	category := params["category"]
	slug  := params["product"]

	if category == "" || slug == "" {
		// if there are no paramters, render 404 error
		renderError(404, r)
	} else {
		// look for the product
		if product, err := app.GetProduct(category, slug); err == nil {
			r.HTML(200, "product", product)
		} else {
			// product not found
			renderError(404, r)
		}
	}
}

// Render error page
func renderError(code int, r render.Render) {
	r.HTML(404, "404", Page{
		HttpCode: 404,
		HttpMethod: "GET",
		Template: "404",
		Path: "/404",
		Slug: "404",
		Title: "Página não encontrada - Hera Modas",
		Description: "A página que você tentou acessar não existe ou foi removida.",
	})
}