# goHermes

GoHermes is a golang written script that scraps scientific documents from several sources. It builds a library for several purposes like learning, AI training, collection, etc.
This tool provides fast document downloading, using Go concurrency, data storage into a Sqlite DB, and a logging system to track downloading errors, like typically withdrawn paper.

# Features

goHermes can currently download pdf from the following sources :
- ePrint : https://eprint.iacr.org/complete/
- free haven : https://www.freehaven.net/anonbib/date.html

# TODO

- Download security : limit the number of concurrent request by domain name
- Download source interface : define an interface for every soruce. The user will then implement its own node parser for each source.
