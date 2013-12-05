package goshop

import (
	"goshop/app/database"
	"labix.org/v2/mgo"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/gzip"
	"github.com/codegangsta/martini-contrib/render"
	"reflect"
	"errors"
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
func (gs *Shop) Start() *Shop {
	// don't connect if is connected
	if gs.databaseSession == nil {
		gs.databaseSession = database.Connect("localhost", "heramodas")
	}

	// create martini instance
	gs.Server = martini.Classic()

	// gzip all requests
	gs.Server.Use(gzip.All())

	return gs
}

// Get category by slug in database
func (gs *Shop) GetCategory(slug string, products bool) (Category, error) {
	var c Category
	var err error = errors.New("Category not found")

	if len(gs.Categories) > 0 {
		// find category in cached category list
		for _, category := range gs.Categories {
			if category.Slug == slug {
				c = category
				err = nil
				break
			}
		}
	}

	// try to find category in database
	if err != nil {
		err = database.Find("categories", database.M{"slug": slug}).One(&c)
	}

	if products == true && err == nil {
		// look for category products
		database.Find("products", database.M{"category": slug}).All(&c.Products)
	}

	return c, err
}

// Get product details in database
func (gs *Shop) GetProduct(category string, slug string) (Product, error) {
	var p Product

	// finc product
	err := database.Find("products", database.M{"category": category, "slug": slug}).One(&p)

	if err == nil {
		p.Category, err = gs.GetCategory(category, false)
	}

	return p, err
}

// Define routes
func (gs *Shop) Route() *Shop {
	// Register routes
	for i := 0; i < len(gs.Routes); i++ {
		page := gs.Routes[i]

		fields := reflect.Indirect(reflect.ValueOf(page))

		switch fields.FieldByName("HttpMethod").String() {
		case "GET":
			switch page.(type) {
				case Page:
					// convert route to page struct
					route := page.(Page)

					gs.Server.Get(route.Path, func(r render.Render) {
						r.HTML(route.HttpCode, route.Template, route)
					})
				case InternalRoute:
					// convert route to InternalRoute struct
					r := page.(InternalRoute)

					gs.Server.Get(r.Path, r.Run)
			}
		}
	}

	return gs
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