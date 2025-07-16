# Source interface

A source is a repository of documents to download. There is a list of the supported sources right now :

| Sources      | Description      | Link |
| ------------- | ------------- | ------------- |
| ePrint archive | The Cryptology ePrint Archive provides rapid access to recent research in cryptology. | https://eprint.iacr.org/ |
| Haven | The Free Haven Project aims to deploy a system for distributed, anonymous, persistent data storage which is robust against attempts by powerful adversaries to find and destroy any stored data.| https://www.freehaven.net/anonbib/date.html |
| HAL | An archive of scientific research papers | https://hal.science/ |


# Interface

The source interface is composed of the following functions :

| Function      | Parameter      | Return value | Description |
| ------------- | ------------- | ------------- | ------------- |
| Init | engine *Engine | error | Gather metadata such as numbers of documents, years. The fetch is specific for each source. |
| Fetch | engine *Engine | error | Start to download every selected documents |

The Init() function is gathering data to setup urls for documents to download.

# Source implementation

For each sources, the Init() and Fetch() function are implemented depending on how the documents are presented. For instance, in the eprint archive paper are organized in years. So in the Init() I'm gathering every years to prepare my Fetch() function.
For those sources there is no public API so I'm scrapping each documents by hand.

This is a typical implementation of a Source :

```go
type FreeHavenSource struct {
	Name           string
	Path           string
	BaseUrl        string
	TotalDocuments int
	Documents      []*Document
}
```
The four mandatory fields are **Name**, **Path**, **BaseUrl** and **Documents**. The name of the source, the path were to store the documents, the baseUrl were to download the documents and finally a list of every documents to Download.
The last field is **TotalDocuments** which is only an informative data.



