package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	v1 "github.com/awwithro/makemea/api/v1"
	"github.com/awwithro/makemea/randomtable"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/justinian/dice"
	"github.com/slack-go/slack"
)

// NewServer returns a gin server that will serve items from the given tree
func NewServer(tree *randomtable.Tree) *gin.Engine {
	e := gin.Default()
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	AttachHandlers(e, tree)
	return e
}

func AttachHandlers(e *gin.Engine, tree *randomtable.Tree) {
	v1 := e.Group("v1")
	v1.GET("/items/*path", getFunc(tree))
	v1.GET("/tables/*path", listFunc(tree))
	v1.GET("/roll/*roll", rollFunc())
	e.POST("/slack/events", slashCommandFunc(tree))
}

func listFunc(tree *randomtable.Tree) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Param("path")
		path = strings.TrimPrefix(path, "/")
		tables := tree.ListTables(path, false)

		c.JSON(http.StatusOK, v1.ListTableResponse{
			Tables: tables,
		})
	}
}

func getFunc(tree *randomtable.Tree) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Param("path")
		path = strings.TrimPrefix(path, "/")
		item, err := tree.GetItem(path)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusOK, v1.GetItemResponse{
			Item: item,
		})
	}
}
func rollFunc() func(*gin.Context) {
	return func(c *gin.Context) {
		roll := c.Param("roll")
		result, _, err := dice.Roll(roll)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, v1.RollResponse{
			Result:      result.Int(),
			Description: result.Description(),
		})
	}
}

func slashCommandFunc(tree *randomtable.Tree) func(*gin.Context) {
	return func(c *gin.Context) {
		newTree := tree.WithStringFormatter()
		tree = &newTree
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
