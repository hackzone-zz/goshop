package goshop

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

// Category router
func (gs *Shop) RouteCategory(params martini.Params, r render.Render) {
	slug := params["category"]

	if slug == "" {
		// if there is no parameter, render 404 error
		renderError(404, r)
	} else {
		// look for the category and their products
		if category, err := gs.GetCategory(slug, true); err == nil {
			// render category page
			r.HTML(200, "category", category)
		} else {
			// category not found
			renderError(404, r)
		}
	}
}

// Product details router
func (gs *Shop) RouteProduct(params martini.Params, r render.Render) {
	category := params["category"]
	slug  := params["product"]

	if category == "" || slug == "" {
		// if there are no paramters, render 404 error
		renderError(404, r)
	} else {
		// look for the product
		if product, err := gs.GetProduct(category, slug); err == nil {
			r.HTML(200, "product", product)
		} else {
			// product not found
			renderError(404, r)
		}
	}
}

