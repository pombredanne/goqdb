package controllers

import (
	"github.com/PacketFire/goqdb/app/models"
	"github.com/PacketFire/goqdb/app/routes"
	"github.com/robfig/revel"
	"strings"
	"time"
)

type App struct {
	GorpController
}

func (c App) Index(page models.PageState) revel.Result {

	var savedAuthor string

	if author, ok := c.Session["author"]; ok {
		savedAuthor = author
	}

	if page.Page == 0 {
		page.Page = 1
	}

	size := page.Size
	if size <= 0 {
		size = 5
	}

	// for pagination
	nextPage := page.Page + 1
	prevPage := page.Page - 1

	page.Search = strings.TrimSpace(page.Search)

	var entries []*models.QdbEntry

	if page.Search == "" {
		entries = loadEntries(c.Txn.Select(models.QdbEntry{},
			`SELECT * FROM QdbEntry ORDER BY QuoteId DESC LIMIT ?, ?`, (page.Page-1)*(size-1), size))
	} else {
		page.Search = strings.ToLower(page.Search)
		entries = loadEntries(c.Txn.Select(models.QdbEntry{},
			`SELECT * FROM QdbEntry WHERE LOWER(Quote) LIKE ? ORDER BY QuoteId DESC LIMIT ?, ?`, "%"+page.Search+"%", (page.Page-1)*(size), size))

	}

	hasPrevPage := page.Page > 1
	hasNextPage := len(entries) == size
	if hasNextPage {
		entries = entries[:len(entries)-1]
	}

	return c.Render(entries, savedAuthor, page, hasPrevPage, prevPage, hasNextPage, nextPage)
}

func (c App) Post(entry models.QdbEntry, page models.PageState) revel.Result {

	entry.Created = time.Now().Unix()
	entry.Rating = 0

	c.Session["author"] = entry.Author

	c.Validation.Required(entry.Quote)
	c.Validation.Required(entry.Author)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Index(page))
	}
	c.Txn.Insert(&entry)
	return c.Redirect(routes.App.Index(page))
}

func (c App) RatingUp(id int, page models.PageState) revel.Result {
	_, err := c.Txn.Exec("UPDATE QdbEntry SET Rating = Rating + 1 WHERE QuoteId = ?", id)

	if err != nil {
	}

	return c.Redirect(routes.App.Index(page))
}

func (c App) RatingDown(id int, page models.PageState) revel.Result {
	_, err := c.Txn.Exec("UPDATE QdbEntry SET Rating = Rating - 1 WHERE QuoteId = ?", id)

	if err != nil {
	}

	return c.Redirect(routes.App.Index(page))
}

func (c App) One(id int) revel.Result {
	var entries []*models.QdbEntry
	entries = loadEntries(c.Txn.Select(models.QdbEntry{},
		`SELECT * FROM QdbEntry WHERE QuoteId = ? ORDER BY QuoteId DESC LIMIT 1`, id))
	if len(entries) == 0 {
		c.Flash.Error("no such id")
	}

	quote := entries[0]

	return Utf8Result(quote.Quote)
}
