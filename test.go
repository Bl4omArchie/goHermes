package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func detectFormatType(path string) string {
		slash := strings.LastIndex(path, "/")
		dot := strings.LastIndex(path, ".")
		if dot > slash && dot != -1 {
				return strings.ToUpper(path[dot+1:])
		}
		return "UNKNOWN"
}

func main() {
	const htmlData = `

	<!DOCTYPE html>
	<html lang="en">
	  <head>
	    

	    <meta charset="utf-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	    <link href="/css/dist/css/bootstrap.min.css" rel="stylesheet">
	    <title>Universally Composable On-Chain Quadratic Voting for Liquid Democracy</title>
	    <link rel="stylesheet" href="/css/eprint.css?v=10">

	    <link rel="shortcut icon" href="/favicon.ico" type="image/x-icon" />
	    <link rel="apple-touch-icon" href="/img/apple-touch-icon-180x180.png" />
	    
	<style>
	  a.toggle-open:after {
	  content:' -';
	  font-weight: 800;
	  }
	  a.toggle-closed:after {
	  content: " ›";
	  font-weight: 800;
	  }
	  .paper-abstract {
	     white-space: pre-wrap;
	  }
	  #metadata dt {
	    margin-top: 1rem;
	  }
	  #metadata dt + dd { /* gap between dt and first dd */
	    margin-top: .75rem;
	  }
	  #metadata dd {
	    margin-left: 2rem;
	  }
	  #metadata dd.keywords {
	    padding-bottom: .5rem;
	  }
	  span.authorName {
	    margin-top: .5rem;
	    font-style: italic;
	  }
	</style>

	<script>
	MathJax = {
	    tex: {
	        inlineMath: [['$', '$'], ['\\(', '\\)']],
	        displayMath: [ ['$$','$$'], ["\\[","\\]"] ],
	        processEnvironments: false
	    },
	    loader: {
	        load: [
	            "ui/safe",
	            "ui/lazy",
	        ],
	    },
	    options: {
	        safeOptions: {
	        allow: {
	            URLs: "none",
	            classes: "safe",
	            cssIDs: "safe",
	            styles: "safe",
	        },
	        },
	    }
	    };
	</script>

	<script id="MathJax-script" async src="/js/mathjax/tex-chtml.js"></script>


	  <meta name="citation_title" content="Universally Composable On-Chain Quadratic Voting for Liquid Democracy">
	  
	  <meta name="citation_author" content="Lyudmila Kovalchuk">
	  
	  <meta name="citation_author" content="Bingsheng Zhang">
	  
	  <meta name="citation_author" content="Andrii Nastenko">
	  
	  <meta name="citation_author" content="Zeyuan Yin">
	  
	  <meta name="citation_author" content="Roman Oliynykov">
	  
	  <meta name="citation_author" content="Mariia Rodinko">
	  
	  <meta name="citation_journal_title" content="Cryptology ePrint Archive">
	  <meta name="citation_publication_date" content="2025">
	  
	  <meta name="citation_pdf_url" content="https://eprint.iacr.org/2025/803.pdf">
	  


	  
	  <meta property="og:image" content="https://eprint.iacr.org/img/iacrlogo.png"/>
	  <meta property="og:image:alt" content="IACR logo"/>
	  <meta property="og:url" content="https://eprint.iacr.org/2025/803">
	  <meta property="og:site_name" content="IACR Cryptology ePrint Archive" />
	  <meta property="og:type" content="article" />
	  
	  <meta property="og:title" content="Universally Composable On-Chain Quadratic Voting for Liquid Democracy" />
	  <meta property="og:description" content="Decentralized governance plays a critical role in blockchain communities, allowing stakeholders to shape the evolution of platforms such as Cardano, Gitcoin, Aragon, and MakerDAO through distributed voting on proposed projects in order to support the most beneficial of them. In this context, numerous voting protocols for decentralized decision-making have been developed, enabling secure and verifiable voting on individual projects (proposals). However, these protocols are not designed to support more advanced models such as quadratic voting (QV), where the voting power, defined as the square root of a voter’s stake, must be distributed among the selected by voter projects. Simply executing multiple instances of a single-choice voting scheme in parallel is insufficient, as it can not enforce correct voting power splitting. To address this, we propose an efficient blockchain-based voting protocol that supports liquid democracy under the QV model, while ensuring voter privacy, fairness and verifiability of the voting results. In our scheme, voters can delegate their votes to trusted representatives (delegates), while having the ability to distribute their voting power across selected projects. We model our protocol in the Universal Composability framework and formally prove its UC-security under the Decisional Diffie–Hellman (DDH) assumption. To evaluate the performance of our protocol, we developed a prototype implementation and conducted performance testing. The results show that the size and processing time of a delegate’s ballot scale linearly with the number of projects, while a voter’s ballot scales linearly with both the number of projects and the number of available delegation options. In a representative setting with 64 voters, 128 delegates and 128 projects, the overall traffic amounts to approximately 2.7 MB per voted project, confirming the practicality of our protocol for modern blockchain-based governance systems." />
	  <meta property="article:section" content="PROTOCOLS" />
	  
	  <meta property="article:modified_time" content="2025-05-05T16:59:57+00:00" />
	  <meta property="article:published_time" content="2025-05-05T16:59:57+00:00" />
	  
	  
	  
	      <meta property="article:tag" content="Quadratic Voting" />
	  
	      <meta property="article:tag" content="Cryptographic Protocol" />
	  
	      <meta property="article:tag" content="Decentralized Decision-making" />
	  
	      <meta property="article:tag" content="Liquid Democracy" />
	  
	      <meta property="article:tag" content="UC-security" />
	  
	      <meta property="article:tag" content="Blockchain" />
	  
	  

	  </head>

	  <body>
	    <noscript>
	      <h1 class="text-center">What a lovely hat</h1>
	      <h4 class="text-center">Is it made out of <a href="https://iacr.org/tinfoil.html">tin foil</a>?</h4>
	    </noscript>
	    <div class="fixed-top" id="topNavbar">
	      <nav class="navbar navbar-custom navbar-expand-lg">
	        <div class="container px-0 justify-content-between justify-content-lg-evenly">
	          <div class="order-0 align-items-center d-flex">
	            <button class="navbar-toggler btnNoOutline" type="button" data-bs-toggle="collapse" data-bs-target="#navbarContent" aria-controls="navbarContent" aria-expanded="false">
	              <span class="icon-bar top-bar"></span>
	              <span class="icon-bar middle-bar"></span>
	              <span class="icon-bar bottom-bar"></span>
	            </button>
	            <a class="d-none me-5 d-lg-inline" href="https://iacr.org/"><img class="iacrlogo" src="/img/iacrlogo_small.png" alt="IACR Logo" style="max-width:6rem;"></a>
	          </div>
	          <a class="ePrintname order-1" href="/">
	            <span class="longNavName">Cryptology ePrint Archive</span>
	          </a>
	          <div class="collapse navbar-collapse order-3" id="navbarContent">
	            <ul class="navbar-nav me-auto ms-2 mb-2 mb-lg-0 justify-content-end w-100">
	              <li class="ps-md-3 nav-item dropdown">
	                <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
	                  Papers
	                </a>
	                <ul class="dropdown-menu me-3" aria-labelledby="navbarDropdown">
	                  <span class="text-dark mx-3" style="white-space:nowrap;">Updates from the last:</span>
	                  <li><a class="dropdown-item ps-custom" href="/days/7">7 days</a></li>
	                  <li><a class="dropdown-item ps-custom" href="/days/31">31 days</a></li>
	                  <li><a class="dropdown-item ps-custom" href="/days/183">6 months</a></li>
	                  <li><a class="dropdown-item ps-custom" href="/days/365">365 days</a></li>
	                  <li><hr class="dropdown-divider"></li>
	                  <li><a class="dropdown-item" href="/byyear">Listing by year</a></li>
	                  <li><a class="dropdown-item" href="/complete">All papers</a></li>
	                  <li><a class="dropdown-item" href="/complete/compact">Compact view</a></li>
	                  <li><a class="dropdown-item" href="https://www.iacr.org/news/subscribe">Subscribe</a></li>
	                  <li><hr class="dropdown-divider"></li>
	                  <li><a class="dropdown-item" href="/citation.html">How to cite</a></li>
	                  <li><hr class="dropdown-divider"></li>
	                  <li><a class="dropdown-item" href="/rss">Harvesting metadata</a></li>
	                </ul>
	              </li>
	              <li class="ps-md-3 nav-item dropdown">
	                <a class="nav-link dropdown-toggle" href="#" id="submissionsDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
	                  Submissions
	                </a>
	                <ul class="dropdown-menu me-3" aria-labelledby="submissionsDropdown">
	                  <li><a class="dropdown-item" href="/submit">Submit a paper</a></li>
	                  <li><a class="dropdown-item" href="/revise">Revise or withdraw a paper</a></li>
	                  <li><a class="dropdown-item" href="/operations.html">Acceptance and publishing conditions</a></li>
	                </ul>
	              </li>
	              <li class="ps-md-3 nav-item dropdown">
	                <a class="nav-link dropdown-toggle" href="#" id="aboutDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
	                  About
	                </a>
	                <ul class="dropdown-menu me-3" aria-labelledby="aboutDropdown">
	                  <li><a class="dropdown-item" href="/about.html">Goals and history</a></li>
	                  <li><a class="dropdown-item" href="/news.html">News</a></li>
	                  <li><a class="dropdown-item" href="/stats">Statistics</a></li>
	                  <li><a class="dropdown-item" href="/contact.html">Contact</a></li>
	                </ul>
	              </li>
	            </ul>
	          </div>
	          <div class="dropdown ps-md-2 text-right order-2 order-lg-last">
	            <button class="btn btnNoOutline" type="button" id="dropdownMenuButton1" data-bs-toggle="dropdown" aria-expanded="false">
	              <img src="/img/search.svg" class="searchIcon" alt="Search Button"/>
	            </button>
	            <div id="searchDd" class="dropdown-menu dropdown-menu-end p-0" aria-labelledby="dropdownMenuButton1">
	              <form action="/search" method="GET">
	                <div class="input-group">
	                  <input id="searchbox" name="q" type="search" class="form-control" autocomplete="off">
	                  <button class="btn btn-secondary border input-group-append ml-2">
	                    Search
	                  </button>
	                </div>
	              </form>
	              <div class="ms-2 p-1 d-none"><a href="/search">Advanced search</a></div>
	            </div>
	          </div>
	        </div>
	      </nav>
	    </div>
	    <main id="eprintContent" class="container px-3 py-4 p-md-4">

	<div class="row mt-4">
	  <div class="col-md-7 col-lg-8 pe-md-5">
	    
	    <h4>Paper 2025/803</h4>
	    <h3 class="mb-3">Universally Composable On-Chain Quadratic Voting for Liquid Democracy</h3>
	    
	    
	    
	    <div class="author"><span class="authorName">Lyudmila Kovalchuk</span><a class="ms-1" target="_blank" href="https://orcid.org/0000-0003-2874-7950"><img class="align-baseline orcidIcon" src="/img/orcid.svg"></a><span class="affiliation">, IOG</span><span class="affiliation">, Pukhov Institute for Modelling in Energy Engineering</span></div>
	    
	    
	    <div class="author"><span class="authorName">Bingsheng Zhang</span><a class="ms-1" target="_blank" href="https://orcid.org/0000-0002-2320-9582"><img class="align-baseline orcidIcon" src="/img/orcid.svg"></a><span class="affiliation">, IOG Research</span><span class="affiliation">, Zhejiang University</span></div>
	    
	    
	    <div class="author"><span class="authorName">Andrii Nastenko</span><a class="ms-1" target="_blank" href="https://orcid.org/0000-0001-7780-5331"><img class="align-baseline orcidIcon" src="/img/orcid.svg"></a><span class="affiliation">, IOG</span><span class="affiliation">, Kharkiv National University of Radio Electronics</span></div>
	    
	    
	    <div class="author"><span class="authorName">Zeyuan Yin</span><a class="ms-1" target="_blank" href="https://orcid.org/0000-0002-3722-9605"><img class="align-baseline orcidIcon" src="/img/orcid.svg"></a><span class="affiliation">, Zhejiang University</span></div>
	    
	    
	    <div class="author"><span class="authorName">Roman Oliynykov</span><a class="ms-1" target="_blank" href="https://orcid.org/0000-0002-3494-0493"><img class="align-baseline orcidIcon" src="/img/orcid.svg"></a><span class="affiliation">, IOG</span><span class="affiliation">, V.N.Karazin Kharkiv National University</span></div>
	    
	    
	    <div class="author"><span class="authorName">Mariia Rodinko</span><a class="ms-1" target="_blank" href="https://orcid.org/0000-0003-4692-9811"><img class="align-baseline orcidIcon" src="/img/orcid.svg"></a><span class="affiliation">, IOG</span><span class="affiliation">, V.N.Karazin Kharkiv National University</span></div>
	    
	    
	    <h5 class="mt-3">Abstract</h5>
	    <p style="white-space: pre-wrap;">Decentralized governance plays a critical role in blockchain communities, allowing stakeholders to shape the evolution of platforms such as Cardano, Gitcoin, Aragon, and MakerDAO through distributed voting on proposed projects in order to support the most beneficial of them. In this context, numerous voting protocols for decentralized decision-making have been developed, enabling secure and verifiable voting on individual projects (proposals). However, these protocols are not designed to support more advanced models such as quadratic voting (QV), where the voting power, defined as the square root of a voter’s stake, must be distributed among the selected by voter projects. Simply executing multiple instances of a single-choice voting scheme in parallel is insufficient, as it can not enforce correct voting power splitting. To address this, we propose an efficient blockchain-based voting protocol that supports liquid democracy under the QV model, while ensuring voter privacy, fairness and verifiability of the voting results. In our scheme, voters can delegate their votes to trusted representatives (delegates), while having the ability to distribute their voting power across selected projects. We model our protocol in the Universal Composability framework and formally prove its UC-security under the Decisional Diffie–Hellman (DDH) assumption. To evaluate the performance of our protocol, we developed a prototype implementation and conducted performance testing. The results show that the size and processing time of a delegate’s ballot scale linearly with the number of projects, while a voter’s ballot scales linearly with both the number of projects and the number of available delegation options. In a representative setting with 64 voters, 128 delegates and 128 projects, the overall traffic amounts to approximately 2.7 MB per voted project, confirming the practicality of our protocol for modern blockchain-based governance systems.</p>
	  
	  </div>
	  <div id="metadata" class="col-md-5 col-lg-4 ps-md-5 mt-4 mt-md-0">
	    <h5>Metadata</h5>
	    <dl>
	      <dt>
	        Available format(s)
	      </dt>
	      <dd>
	        
	        
	        <a class="btn btn-sm btn-outline-dark"
	           href="/2025/803.pdf">
	          <img class="icon" src="/img/file-pdf.svg">PDF</a>
	        
	        
	        
	      </dd>
	  
	      <dt>Category</dt>
	        <dd><a href="/search?category=PROTOCOLS"><small class="badge category category-PROTOCOLS">Cryptographic protocols</small></a></dd>
	  
	      <dt>Publication info</dt>
	      <dd>Preprint. </dd>
	  
	  
	      <dt>Keywords</dt>
	      <dd class="keywords"><a href="/search?q=Quadratic%20Voting" class="me-2 badge bg-secondary keyword">Quadratic Voting</a><a href="/search?q=Cryptographic%20Protocol" class="me-2 badge bg-secondary keyword">Cryptographic Protocol</a><a href="/search?q=Decentralized%20Decision-making" class="me-2 badge bg-secondary keyword">Decentralized Decision-making</a><a href="/search?q=Liquid%20Democracy" class="me-2 badge bg-secondary keyword">Liquid Democracy</a><a href="/search?q=UC-security" class="me-2 badge bg-secondary keyword">UC-security</a><a href="/search?q=Blockchain" class="me-2 badge bg-secondary keyword">Blockchain</a></dd>
	  
	      <dt>Contact author(s)</dt>
	      <dd><span class="font-monospace">
	        lyudmila kovalchuk<span class="obfuscate"> &#64; </span>iohk io<br>bingsheng zhang<span class="obfuscate"> &#64; </span>iohk io<br>andrii nastenko<span class="obfuscate"> &#64; </span>iohk io<br>zeyuanyin<span class="obfuscate"> &#64; </span>zju edu cn<br>roman oliynykov<span class="obfuscate"> &#64; </span>iohk io<br>mariia rodinko<span class="obfuscate"> &#64; </span>iohk io
	        </span>
	      </dd>
	      <dt>History</dt>
	      
	      
	      <dd>2025-05-05: approved</dd>
	      
	      
	      
	      <dd>2025-05-05: received</dd>
	      
	      
	      
	      <dd><a rel="nofollow" href="/archive/versions/2025/803">See all versions</a></dd>
	      
	      <dt>Short URL</dt>
	      <dd><a href="https://ia.cr/2025/803">https://ia.cr/2025/803</a></dd>
	      <dt>License</dt>
	      <dd><a rel="license" target="_blank" href="https://creativecommons.org/licenses/by/4.0/">
	        <img class="licenseImg" src="/img/license/CC_BY.svg" alt="Creative Commons Attribution" title="Creative Commons Attribution"><br>
	        <small>CC BY</small>
	        </a>
	      </dd>
	    </dl>
	  </div>
	  </div>
	  
	  
	  <p class="mt-4"><strong>BibTeX</strong> <button id="bibcopy" class="ms-2 btn btn-sm btn-outline-dark"
	 aria-label="Copy to clipboard"
	                                     onclick="copyBibtex()">
	    <img src="/img/copy-outline.svg" class="icon">Copy to clipboard</button></p>
	  <pre id="bibtex">
	@misc{cryptoeprint:2025/803,
	      author = {Lyudmila Kovalchuk and Bingsheng Zhang and Andrii Nastenko and Zeyuan Yin and Roman Oliynykov and Mariia Rodinko},
	      title = {Universally Composable On-Chain Quadratic Voting for Liquid Democracy},
	      howpublished = {Cryptology {ePrint} Archive, Paper 2025/803},
	      year = {2025},
	      url = {https://eprint.iacr.org/2025/803}
	}
	</pre>
	  

	<script>
	 var bibcopy;
	 function triggerTooltip() {
	   console.log('setting tooltip');
	 }
	 window.onload = triggerTooltip;
	 function copyBibtex() {
	   let range = document.createRange();
	   range.selectNode(document.getElementById('bibtex'));
	   window.getSelection().removeAllRanges();
	   window.getSelection().addRange(range);
	   document.execCommand('copy');
	   window.getSelection().removeAllRanges();
	   let bibcopy = document.getElementById('bibcopy');
	   let copyTooltip = new bootstrap.Tooltip(bibcopy, {trigger: 'manual',
	                                                     title: 'Copied!'});
	   copyTooltip.show();
	   setTimeout(function() {
	     copyTooltip.dispose();
	   }, 2000);
	 }
	</script>

	</main>
	<div class="container-fluid mt-auto" id="eprintFooter">
	  <a href="https://iacr.org/">
	    <img id="iacrlogo" src="/img/iacrlogo_small.png" class="img-fluid d-block mx-auto" alt="IACR Logo">
	  </a>
	  <div class="colorDiv"></div>
	  <div class="alert alert-success w-75 mx-auto">
	    Note: In order to protect the privacy of readers, eprint.iacr.org
	    does not use cookies or embedded third party content.
	  </div>
	</div>
	<script src="/css/bootstrap/js/bootstrap.bundle.min.js"></script>
	  <script>
	   var topNavbar = document.getElementById('topNavbar');
	   if (topNavbar) {
	     document.addEventListener('scroll', function(e) {
	       if (window.scrollY > 100) {
	         topNavbar.classList.add('scrolled');
	       } else {
	         topNavbar.classList.remove('scrolled');
	       }
	     })
	   }
	  </script>


	</body>
	</html>
	`

		doc, err := html.Parse(strings.NewReader(htmlData))
		if err != nil {
				log.Fatal(err)
		}

		formats := make(map[string]string)

		var inMetadata, inFormatSection bool

		var f func(*html.Node)
		f = func(n *html.Node) {
				if n.Type == html.ElementNode {
					// Look for the Metadata <div>
					if n.Data == "div" {
						for _, attr := range n.Attr {
							if attr.Key == "id" && attr.Val == "metadata" {
								inMetadata = true
								break
							}
						}
					}

					if inMetadata && n.Data == "dt" {
						if n.FirstChild != nil && strings.Contains(strings.ToLower(n.FirstChild.Data), "available format") {
							inFormatSection = true
						}
					}

					if inFormatSection && n.Data == "a" {
						href := ""
						for _, attr := range n.Attr {
							if attr.Key == "href" {
								href = attr.Val
							}
						}
						if href != "" {
							// Get file extension or type
							formatType := detectFormatType(href)
							formats[formatType] = href
						}
					}
				}

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
		}
		f(doc)

		baseURL := "https://eprint.iacr.org/2025/803"
		for format, path := range formats {
				fmt.Printf("%s: %s\n", format, baseURL+path)
		}
}
