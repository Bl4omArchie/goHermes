# ePrint Database Project

Goal : getting every ePrint papers in a database the fastest as possible.

Technologies :
- scripting : golang
- SGBD : postgreSQL
- concurrency

Result of retrieving datas from every 2024 papers :
```
Execution time: 4m43.4106983s
```

Result of retrieving datas from every 2024, 2023 and 2022 papers in concurrency :
```
Total execution time = ~10m40s
```
For this test, I launched three goroutines one for each year. I think it is a poor concurrency design and I can find better.

Notes : 
- too much goroutines making simultaneous request at one server lead to this error:
read: connection reset by peer

- Find a better design for my concurrency