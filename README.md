goqdb
=====

Quote Database in Go

usage
=====

requirements
------------
* http://robfig.github.io/revel/
* https://github.com/mattn/go-sqlite3
* working sqlite3 library

running
-------

```
revel run github.com/PacketFire/goqdb
```

api
---

### QdbEntry ###

<table>
	<thead>
		<tr>
			<th>Name</th> <th>Type</th> <th>Description</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>QuoteId</td> <td>int</td> <td>The quote id</td>
		</tr>
		<tr>
			<td>Quote</td> <td>string</td> <td>The quote body</td>
		</tr>
		<tr>
			<td>Author</td> <td>string</td> <td>The author of the quote</td>
		</tr>
		<tr>
			<td>Created</td> <td>int64</td> <td>unix time in seconds</td>
		</tr>
		<tr>
			<td>Rating</td> <td>int</td> <td>The quote's rating</td>
		</tr>
	</tbody>
</table>

All resources return *200* on success or *500* with an undefined body 
if fatal errors were encountered. Resources requiring an id return a 
*404* with undefined body if the id does not exist in the database. 

### Retrieve the entire database
	
	GET /api/v0

### Insert a new entry

	POST /api/v0

Accepts *Quote* and *Author* fields of a *QdbEntry*

Note: POST returns 201 Created on success and 400 Bad Request
if the post data did not pass validation

### Retrieve quote entry

*:id* is used here in place of the quote id for the target entry

	GET /api/v0/:id

### Upvote a quote

	PUT /api/v0/:id/rating

### Downvote a quote

	DELETE /api/v0/:id/rating

