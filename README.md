# ePrint D.

Eprint DB is a golang tool that scrap and download the Eprint archive. My intention is to gather many data as possible for my personal uses and projects. Its also a practise for Golang language.

Current the script support the following sources :
- ePrint

And there is the in-coming sources :
- arxiv
- NIST
- HAL

# Roadmap

- [x] Scrapping system for Eprint
- [x] Log system
- [x] Worker pool download
- Ingest data into sqliteDB
    - [x] Documents
    - [] License (issue)
    - [] Authors
- [] Fix folder and file creation issue
- [] Make the worker pool stop
- [] Add more sources like arxiv, NIST...

# Preview of the tool interface

!["picture"](img/menu.png)