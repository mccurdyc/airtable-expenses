# airtable-expenses
---

This is a project that I created to test out [Airtable](https://airtable.com/),
a spreadsheet-like UI, with database properties such as relationships between
sheets. This specific project is the code that I used to migrate from a Google Sheet
with a single sheet to Airtable with multiple tables.

## Important Notes

[Here is a link](https://airtable.com/shrRW6c3cMnoiz3KB) to the Airtable Base that I created
and can be copied.

Airtable generates an API specific to your table, see the documentation [here](https://airtable.com/api).
> After you’ve created and configured the schema of an Airtable base from the graphical
> interface, your Airtable base will provide its own API to create, read, update,
> and destroy records.

With that said, use this as an example as it is not a generalized tool.

## Dependencies
+ [Go](https://golang.org/doc/install)
+ Obtain the Google libraries used
  ```bash
  go get -u google.golang.org/api/sheets/v4
  go get -u golang.org/x/oauth2/...
  ```

## References
+ [Google Sheets Go Quickstart](https://developers.google.com/sheets/api/quickstart/go)

## Getting Started
1. Click the "ENABLE THE SHEETS API" button in the [Google Quickstart Guide]((https://developers.google.com/sheets/api/quickstart/go))
2. "Clone" the project
  ```bash
  go get -u github.com/mccurdyc/cmd/migrate
  ```
3. Setup the environment
  1. Fill in the values in `.env.sample`
  2. `source .env.sample`
4. Run the project
  ```bash
  make run
  ```

## Comments on Airtable After Usage

### What I like

### What I dislike
+ can only GET by Airtable ID, which means that I have to GET ALL and store ID in memory
+ no `UNIQUE` fields like in a relational database
+ some JSON fields are camelCase, others PascalCase in the API
+ in my experience, in the iOS mobile application, you can only search records
  by the primary key, which I have set as an auto-incremented integer instead of a string.
  + for this, I have considered setting the primary keys as the names of merchants, tags, etc.,
    but I am used to IDs as the primary key.

## License
+ [GNU General Public License Version 3](./LICENSE)
