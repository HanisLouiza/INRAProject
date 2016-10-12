/*
	Package catalog
	Extract names of genes defined in catalog file
*/
package catalog

import (
	"bufio"
	"os"
	"strings"
	"../checkError"
)

func ExtractGenesNames(genePath string) []string {

	if genePath == "" {
		genePath = "./data/genes.txt"
	}
	InputfileGene, err := os.Open(genePath)
	checkError.Check(err)
	defer InputfileGene.Close()

	readerGeneFile := bufio.NewReader(InputfileGene)
	scannerGeneFile := bufio.NewScanner(readerGeneFile)

	genesNames := make([]string, 0)

	for scannerGeneFile.Scan() {
		Line := scannerGeneFile.Text()
		if len(Line) > 0 {
			Words := strings.Fields(Line)
			if len(Words) > 0 {
				genesNames = append(genesNames, Words[1])
			}
		}
	}
	return genesNames
}
