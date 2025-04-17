# ePrint D.

This tool is a golang script that download papers from the website ePrint.
It stores every PDF into a datalake and apply compression for a better storage.

# TODO
In order of priority :

1- Concurrency downloading ✅
    1.1- DownloadPool structure ✅
    1.2- Goroutines for download ✅
    1.3- Limit rate to avoid timeout ✅

2- Data parsing 🚧
    2.1- Retrieve authors names
    2.2- Title of the document
    2.3- Category
    2.4- Release date

3- Database
    3.1- Creation of database with Sqlite (script and schema) 🟠
    3.2- Start filling for testing 🟠
    3.3- begin workflow with filling db 🚧

4- Error channel 🚧
    3.1- reporting missing documents
    3.2- handle error to continue the download


Bonus :
- Progress bar ✅


# Workflow

harvest documents -> store the documents -> ingest documents in DB -> update datalake (compression)