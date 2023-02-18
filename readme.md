# dbstream

Simple library for creating streams from sql query result sets.

#installation

  Note: This library uses golang generics so works for golng 1.18 and above.

    go get github.com/cybernostics/dbstream

# usage

See the stream_test.go file for an example using an sql.db connection to an sqlite database.
The pattern is the same regardless of the type of database.

Once you have executed an sql.db query you call the following method to get the stream:

    rowStream := dbstream.RowStreamFrom(rows, scanFunction)
	
The scanFunction is a function that reads a line from the database using the rows.Scan function

Given a Person stuct the scanPerson function looks like this:

```golang
func scanPerson(rows *sql.Rows, p *Person) error {
	return rows.Scan(&p.Id, &p.FirstName, &p.LastName, &p.Email, &p.IpAddress)
}
```

The Stream object has the same looping behavior as the sql.Rows object. ie call next() on each row and then call Current() to get the result if Next returns true.

Also a call to Close ensures that the rows object is closed.

# background

Why use streams at all?
The stream library encourages users to encapsulate the function for extracting row results to a structure. By putting this in one place, changes to the structure only need to be made in one place.

Also, by using streams instead of collections, you can avoid the need to allocate large amounts of memory by only working with a single row at a time.