# Saloon

[![Build Status](https://travis-ci.org/go-saloon/saloon.svg?branch=master)](https://travis-ci.org/go-saloon/saloon)
[![GoDoc](https://godoc.org/github.com/go-saloon/saloon?status.svg)](https://godoc.org/github.com/go-saloon/saloon)

Saloon is a _Work in Progress_ forum based on [Buffalo](https://gobuffalo.io).

## Database setup

```
$> docker run --name saloon-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```

### Create Your Databases

Ok, so you've edited the "database.yml" file and started postgres, now Buffalo can create the databases in that file for you:

```
$> buffalo db create -a
```

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

```
$> buffalo dev
```

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.
