# ePrint Database Project

Goal : getting every ePrint papers in a database the fastest and the cleaner as possible.
Later goal : have a nice web UI where you can download PDF from several sources (eprint, NIST...), make refine tuning with ollama or get alerted from new papers.

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

# Design ideas

Stages :
1) Retrieve data about a paper : title, pdf url, category
2) Download the PDF
3) store the raw binary in the database

Draw of the process :
```

Start: GetPapersYear -> For each papers: RetrieveDataPaper -> DownloadPaper -> InsertBinary (into the dabase)

```

First idea :
- A fixed number N of goroutines for stages 1, 2 and 3 <br/>
(ie : I create 100 goroutines for each stages and when they have finished one task, they continue with the next one)

Second idea :
- N goroutines for each stages that cannot exceed a limit P of goroutines. <br/>
(ie : I create goroutines for my task until I reached the limit P. Then I wait that some of them has done to create a new ones) 

More ideas : <br/>
Basically the same but using a pipeline with channels. A fourth idea idea could be to use custom rating limit with work-stealing queue.
I shall explore and test those ideas.


# Statistics 

## Introduction 

Statistics are the data I need to correctly download the precise amount of papers available on ePrint. For now the code is a bit goofy because I'm doing a request each time the tool is executed for retrieving data such as categories names, number of PDF from past years. Futur improving will be to have a fine strategy for optimizing the number of request and winning time by calling intelligently the website.


## Rate limit issue

In order to anticipate rate limit issue (from hardware or the ePrint server), I decided to make a small analysis of how many request I will need, how many insertion in my database etc

There is the volume of paper for each years :
```
"2024":1799, "2023":1971, "2022":1781, "2021":1705, "2020":1620,
"2019":1498, "2018":1249, "2017":1262, "2016":1195, "2015":1255, "2014":1029, "2013":881, "2012":733, "2011":714, "2010":660, 
"2009":638, "2008":545, "2007":482, "2006":485, "2005":469, "2004":375, "2003":265, "2002":195, "2001":113, "2000":69,
"1999":24, "1998":26, "1997":15, "1996": 16,

```

Years between 2014 and 2024 have more than one thousand papers which I consider to be the years who need more goroutines.
In the other case years, between 1996 and 2013, there is only less than one thousand papers or even a few dozen which means we don't need too much goroutines.


# Alerts

I'm developing an alert system that allows through a channel to communicate errors and failures. For the moment the system is very simple and I still need to make improvements but I'm working hard on designing a good system that handle every situation.

Currently, my system is a predermined set of flags but it is a bit naive. Normally every cased are handle but what if I get an unknow error ? For instance, while I was making my test for downloading with goroutines, I got a rejection from ePrint website. Now I know this error exists but it means I need to anticipate every possible error and also the case where idk what it is.

A second point is the exit action. With my alerts system you can chose if you want to continue your program or quit. This is okay for the moment, but when my program will run with hundred of gouroutines, how am I going to manage the exiting of all those goroutines ?

Finnaly, the third point is about the strategy behind continuing the program even after a failed attempt of downloading a PDF. While I was correcting my code I saw that I was continuing my script in cases where I already knew the url was incorrect. Like a switch case but without break. So the program was running on the same url for nothing. I need to find better way to continue my program and skipping immediatly when an error about wrong url occurs. And even for the rejected connection I was talking about latly: How can I stop temporaly my program, keep thing frozed where I was, wait a bit and then continue like nothing happened ? 

# Docker image

My docker image is still under development. You can see it as a test.


# Sources 
- https://medium.com/novai-go-programming-101/running-a-golang-application-with-docker-and-docker-compose-2e8d6ab41bde
- https://medium.com/@jamal.kaksouri/the-complete-guide-to-context-in-golang-efficient-concurrency-management-43d722f6eaea
- https://snyk.io/fr/blog/containerizing-go-applications-with-docker/
- https://medium.com/novai-go-programming-101/running-a-golang-application-with-docker-and-docker-compose-2e8d6ab41bde
- https://dev.to/francescoxx/build-a-crud-rest-api-in-go-using-mux-postgres-docker-and-docker-compose-2a75