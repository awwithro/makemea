package server

import (
	"net/http"
	"strings"

	v1 "github.com/awwithro/makemea/api/v1"
	"github.com/awwithro/makemea/randomtable"
	"github.com/gin-gonic/gin"
)

// NewServer returns a gin server that will serve items from the given tree
func NewServer(tree randomtable.Tree) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("v1")
	v1.GET("/items/*path", getFunc(tree))
	v1.GET("/tables/*path", listFunc(tree))
	return r
}

func listFunc(tree randomtable.Tree) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Param("path")
		path = strings.TrimPrefix(path, "/")
		tables := tree.ListTables(path, false)

		c.JSON(http.StatusOK, v1.ListTableResponse{
			Tables: tables,
		})
	}
}

func getFunc(tree randomtable.Tree) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Param("path")
		path = strings.TrimPrefix(path, "/")
		item, err := tree.GetItem(path)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}
		c.JSON(http.StatusOK, v1.GetItemResponse{
			Item: item,
		})
	}
}
