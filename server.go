package main

import (
	"net/http"
	"goshop/app"
	"github.com/codegangsta/martini-contrib/render"
	"html/template"
	//"labix.org/v2/mgo/bson"
	"reflect"
	"log"
)

func main() {
	log.Println("Starting...")

	// create shop app instance
	shop := app.Shop{}

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

	shop.Routes = append(shop.Routes, app.InternalRoute{
		HttpMethod: "GET",
		Path: "/:category",
		Run: shop.RouteCategory,
	})

	shop.Routes = append(shop.Routes, app.InternalRoute{
		HttpMethod: "GET",
		Path: "/:category/:product",
		Run: shop.RouteProduct,
	})

	shop.Start().Route()


	// get all categories
	//database.Find("categories", database.M{}).All(&shop.Categories)

	// set template path
	shop.Server.Use(render.Renderer("templates", func(t *template.Template) {
		t.Funcs(template.FuncMap{
			// Check if the Slug content it's equals some path
			"categorySlug": func(obj interface{}, path string) bool {
				var slug string = ""

				o := reflect.ValueOf(obj)
				fields := reflect.Indirect(o)

				// Try to get the Category struct. It's an exception
				category := fields.FieldByName("Category")

				if category.IsValid() {
					// get slug category
					slug = category.FieldByName("Slug").String()
				} else {
					// it's not a Category struct
					// so, try to get Slug at the first level struct
					slug = fields.FieldByName("Slug").String()
				}

				return slug == path
			}})
	}))

	http.ListenAndServe(":8080", shop.Server)
}
