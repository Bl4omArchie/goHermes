# ePrint Database Project

Goal : getting every ePrint papers in a database

Technologies :
- scripting : golang
- SGBD : postgreSQL
- concurrency

Here is a schema representing my goal for downloading in concurrency each pdf :
```
	     download paper 
start -> download paper -> store the paper
	     download paper 
```