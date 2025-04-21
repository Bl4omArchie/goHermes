package core

import (
	"fmt"
	"github.com/pdfcpu"
	"github.com/pdfcpu"
)


func ReadPDF(filePath string) {
	ctx, err := api.ReadContextFile(filePath, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error reading PDF:", err)
		return
	}
  
	for i := 1; i <= ctx.PageCount; i++ {
		page, err := ctx.ExtractPageText(i)
		if err != nil {
			fmt.Printf("Error extracting text from page %d: %s\n", i, err)
		} else {
			fmt.Printf("Page %d:\n%s\n", i, page)
		}
	}
}