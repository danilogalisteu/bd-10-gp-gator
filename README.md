# bd-10-gp-gator

Guided project:
> Build a blog aggregator microservice in Go. Put your API, database, and web scraping skills to the test.

# gator

This console application allows users to:
* Add RSS feeds from across the internet to be collected
* Store the collected posts in a PostgreSQL database
* Follow and unfollow RSS feeds that other users have added
* View summaries of the aggregated posts in the terminal, with a link to the full post

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## Installation

Clone this repository with:

```bash
git clone https://github.com/danilogalisteu/bd-10-gp-gator.git
```

The application uses the Go language and a local PostgreSQL database.
Following the [instructions here](https://go.dev/doc/install), install the specific Go language runtime for your platform.
Following the [instructions here](https://www.postgresql.org/download/), install a recent version of the specific PostgreSQL server for your platform.

After installing the requirements, start the application with:

```bash
go run .
```

To make a self-contained executable file called `gator`, run:

```bash
go build -o gator
```

To install the application on your system, to be used from any path as the command `gator`, run:

```bash
go install -o gator
```

The required external dependencies should be downloaded and installed automatically on the first run or build.

## Configuration

### Database

The local PostgreSQL server should be run as a service, under a specific user such as `postgres`, with a password.
The Linux commands to set the user password and access the server shell are:
```bash
sudo passwd postgres
sudo service postgresql start
sudo -u postgres psql
```

In the psql shell, under the chosen user, you should create a database named `gator`, usig the SQL query
```sql
CREATE DATABASE gator;
```
and then connect to the database using the command `\c gator`.

Finally, the following command sets the database user and password for connection:
```sql
ALTER USER postgres PASSWORD 'postgres';
```
The database user (referred to as PG_USER) should be the same as the system user, and the database password (referred to as PG_PASS) should be different from the system password.

### Application

The application expects a text file named `.gatorconfig.json` in your home directory.
The file uses the JSON format and should be created manually with the following content:
```json
{
    "db_url":"postgres://<PG_USER>:<PG_PASS>@localhost:5432/gator?sslmode=disable"
}
```

## Use

The following commands are available in the application:

* `gator register USERNAME`: add user
* `gator reset`: remove all users
* `gator login USERNAME`: log in existing user
* `gator users`: list registered users
* `gator agg DURATION`: refresh feeds periodically
* `gator addfeed NAME URL`: add new feed and follow (under logged user)
* `gator feeds`: list saved feeds
* `gator follow URL`: follow existing feed (under logged user)
* `gator following`: list followed feeds (under logged user)
* `gator unfollow URL`: unfollow existing feed (under logged user)
* `gator browse [LIMIT]`: list posts from user feeds (under logged user)
