package controllers

import (
	"github.com/PacketFire/goqdb/app/models"
	"time"
)

var DEFAULT_SIZE = 5

type Core struct {
	GorpController
}

func loadEntries (result []interface{}) []*models.QdbEntry {

	var entries []*models.QdbEntry

	for _, r := range result {
		entries = append(entries, r.(*models.QdbEntry))
	}

	return entries
}

func (c *Core) insertEntry (entry *models.QdbEntry) error {

	entry.Created = time.Now().Unix()
	entry.Rating = 0

	c.Session["author"] = entry.Author

	return c.Txn.Insert(entry)
}

func (c *Core) getEntryById (id int) ([]*models.QdbEntry, error) {

	result, err := c.Txn.Select(models.QdbEntry{},
		`SELECT
			*
		FROM
			QdbEntry
		WHERE
			QuoteId = ?
		LIMIT 1`,
		id)

	return loadEntries(result), err
}

func (c *Core) getEntries (page, size int, search string) ([]*models.QdbEntry, error) {

	if size == 0 {
		size = DEFAULT_SIZE
	}

	var lower int
	if size > 0 {
		lower = size * (page - 1)
	} else {
		lower = 0
	}

	var result []interface{}
	var err error

	if search == "" {
		result, err = c.Txn.Select(models.QdbEntry{},
			`SELECT 
				* 
			FROM 
				QdbEntry 
			ORDER BY 
				QuoteId DESC 
			LIMIT ?, ?`,
			lower, size)
	} else {
		result, err = c.Txn.Select(models.QdbEntry{},
			`SELECT
				*
			FROM
				QdbEntry
			WHERE
				LOWER(Quote)
			LIKE
				?
			LIMIT
				?, ?`,
			"%"+search+"%", lower, size)
	}

	return loadEntries(result), err
}

func (c *Core) upVote (id int) (int64, error) {

	_, err := c.Txn.Exec(
		`UPDATE
			QdbEntry
		SET
			Rating = Rating + 1
		WHERE
			QuoteId = ?`,
		id)

	if err != nil {
		return 0, err
	}

	var changes int64
	changes, err = c.Txn.SelectInt(`SELECT CHANGES()`)

	return changes, err
}

func (c *Core) downVote (id int) (int64, error) {

	_, err := c.Txn.Exec(
		`UPDATE
			QdbEntry
		SET
			Rating = Rating - 1
		WHERE
			QuoteId = ?`,
		id)

	if err != nil {
		return 0, err
	}

	var changes int64
	changes, err = c.Txn.SelectInt(`SELECT CHANGES()`)

	return changes, err
}

