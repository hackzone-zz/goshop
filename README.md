goshop
======

Simple e-commerce system used to learn Go Programming Language

## Dependencies

You need to run `go get` to the following packages:

	labix.org/v2/mgo
	github.com/codegangsta/martini
	github.com/codegangsta/martini-contrib/gzip
	github.com/codegangsta/martini-contrib/render

## Database

I am using MongoDB as Database. This is the database structure:

    collection: categories
    {
    	"slug"        : <string> Unique Category path,
    	"title"       : <string> Category Title,
    	"description" : <string> Category Description
    }

    collection: products
    {
    	"category"    : <string> Unique category path,
    	"slug"        : <string> Unique product path,
    	"title"       : <string> Product path,
    	"description" : <string> Product description,
    	"image"       : <string> Product image path,
    	"price"       : <Float> Product price,
    	"off"         : <Float> Product off
    }
