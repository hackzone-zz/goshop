package main

import (
	"net/http"
	"goshop/app"
	"goshop/app/database"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/gzip"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
	//"labix.org/v2/mgo/bson"
	"reflect"
	"log"
)

/*
// Get on database all categories
func (app *App) getCategories(force bool) []Category {
	if force == true {
		categoriesCollection := app.session.DB("heramodas").C("categories")

		// remove all categories before
		app.categories = app.categories[:0]

		// now, get all categories on database
		categoriesCollection.Find(bson.M{}).All(&app.categories)
	}

	return app.categories
}

// Get the category Slug data
func (app *App) getCategory(slug string) (Category, bool) {
	for _, category := range app.categories {
		if category.Slug == slug {
			return category, true
		}
	}

	return Category{}, false
}

// Get one product on Database
func (app *App) getProduct(category, slug string) (Product, bool) {
	product := Product{}
	collection := app.session.DB("heramodas").C("products")

	collection.Find(bson.M{"slug": slug, "category": category}).One(&product)

	if product.Slug != "" {
		return product, true
	} else {
		return product, false
	}
}

// Get on database all products on category
func (app *App) getCategoryProducts(category *Category) {
	productsCollection := app.session.DB("heramodas").C("products")
	productsCollection.Find(bson.M{"category": category.Slug}).All(&category.Products)
}
*/

func main() {
	log.Println("Starting...")

	// create shop app instance
	shop := app.Shop{
		DatabaseSession: database.Connect("localhost", "heramodas"),
	}

	// get all categories
	database.Find("categories", database.M{}).All(&shop.Categories)

	m := martini.Classic()

	// gzip all requests
	m.Use(gzip.All())

	// set template path
	m.Use(render.Renderer("templates", func(t *template.Template) {
		t.Funcs(template.FuncMap{
			// Check if the Slug content it's equals some path
			"categorySlug": func(obj interface{}, path string) bool {
				var slug string = ""

				o := reflect.ValueOf(obj)
				points := reflect.Indirect(o)

				// Try to get the Category struct. It's an exception
				category := points.FieldByName("Category")

				if category.IsValid() {
					// get slug category
					slug = category.FieldByName("Slug").String()
				} else {
					// it's not a Category struct
					// so, try to get Slug at the first level struct
					slug = points.FieldByName("Slug").String()
				}

				return slug == path
			}})
	}))

	// set static path
	m.Use(martini.Static("assets"))

	shop.Routes = append(shop.Routes, app.Page{
		HttpCode: 200,
		HttpMethod: "GET",
		Template: "page",
		Path: "/",
		Slug: "",
		Title: "Hera Modas",
		Description: "Vendas on-line de vestidos de festa, longos e longuetes, além de bolsas e bijouterias em strass. Trabalhamos também com tamanhos grandes, até XXG.",
	})

	shop.Routes = append(shop.Routes, app.Page{
		HttpCode: 200,
		HttpMethod: "GET",
		Template: "localizacao",
		Path: "/localizacao",
		Slug: "localizacao",
		Title: "Lojas - Hera Modas",
		Description: "Endereço e horários de funcionamento da loja física da Hera Modas e Presentes.",
	})

	shop.Routes = append(shop.Routes, app.Page{
		HttpCode: 404,
		HttpMethod: "GET",
		Template: "404",
		Path: "/404",
		Slug: "404",
		Title: "Página não encontrada - Hera Modas",
		Description: "A página que você tentou acessar não existe ou foi removida.",
	})

	// Register routes
	for i := 0; i < len(shop.Routes); i++ {
		page := shop.Routes[i]

		switch page.HttpMethod {
		case "GET":
			m.Get(page.Path, func(r render.Render) {
				log.Println( page.Template )
				r.HTML(page.HttpCode, page.Template, page)
			})
		}
	}

	/*
	// category
	m.Get("/:category", func(params martini.Params, r render.Render) {
		if category, ok := app.getCategory(params["category"]); ok {
			app.getCategoryProducts(&category)

			r.HTML(200, "category", category)
		} else {
			r.HTML(404, "404", page404)
		}
	})

	// product
	m.Get("/:category/:product", func(params martini.Params, r render.Render) {
		if product, ok := app.getProduct(params["category"], params["product"]); ok {
			product.Category, _ = app.getCategory(params["category"])
			r.HTML(200, "product", product)
		} else {
			r.HTML(404, "404", page404)
		}
	})
	*/


	http.ListenAndServe(":8080", m)
}
