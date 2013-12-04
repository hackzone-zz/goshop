package app

import (
	"labix.org/v2/mgo"
)

type Shop struct {
	DatabaseSession *mgo.Session
	Categories []Category
	Routes []Page
}

type Page struct {
	Title, Description string
	Path, Template, Slug string
	HttpCode int
	HttpMethod string
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


