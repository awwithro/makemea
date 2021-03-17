package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	v1 "github.com/awwithro/makemea/api/v1"
	"github.com/awwithro/makemea/randomtable"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

// NewServer returns a gin server that will serve items from the given tree
func NewServer(tree randomtable.Tree) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("v1")
	v1.GET("/items/*path", getFunc(tree))
	v1.GET("/tables/*path", listFunc(tree))
	r.POST("/slack/events", slashCommandFunc(tree))
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

func slashCommandFunc(tree randomtable.Tree) func(*gin.Context) {
	return func(c *gin.Context) {
		s, err := slack.SlashCommandParse(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, SlackResponse{
				Text:         err.Error(),
				ResponseType: Ephemeral,
			})
			return
		}

		if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		switch s.Command {
		case "/makemea":
			params := &slack.Msg{Text: s.Text}
			words := strings.Split(params.Text, " ")
			if words[0] == "" {
				c.JSON(http.StatusOK, SlackResponse{
					Text:         "Usage:\n 'list [prefix]' to see table names.\n'<table name>' to roll on a table",
					ResponseType: Ephemeral,
				})
				return
			}
			if words[0] == "list" {
				tableName := ""
				if len(words) > 1 {
					tableName = words[1]
				}
				tables := tree.ListTables(tableName, false)
				// Wrap items in back ticks for easier copy/pasting
				for i, t := range tables {
					tables[i] = "`" + t + "`"
				}
				tablesStr := strings.Join(tables, "\n")
				resp := SlackResponse{
					Text:         tablesStr,
					ResponseType: InChannel,
				}
				c.JSON(http.StatusOK, resp)
				return
			} else {
				item, err := tree.GetItem(words[0])
				if err != nil {
					c.JSON(http.StatusOK, SlackResponse{
						Text:         err.Error(),
						ResponseType: Ephemeral,
					})
					return
				}
				c.JSON(http.StatusOK, SlackResponse{
					Text:         fmt.Sprintf("```\n%s\n```", item),
					ResponseType: InChannel,
				})
			}
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
