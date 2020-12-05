package helpers

import "strconv"

type PaginationLink struct {
	Link string
	Name string
}

func GeneratePagination(pgNum int, maxPages int) []PaginationLink {
	pag := make([]PaginationLink, 0)
	var n string
	// delta := 2

	// Add pagLink to prev btn
	if pgNum == 1 {
		pag = append(pag, PaginationLink{Link: "", Name: "←"})
	} else {
		n = strconv.Itoa(pgNum - 1)
		pag = append(pag, PaginationLink{Link: n, Name: "←"})
	}

	if maxPages < 6 {
		for i := 1; i <= maxPages; i++ {
			n = strconv.Itoa(i)
			pag = append(pag, PaginationLink{Link: n, Name: n})
		}
	} else {

	}

	// Add pagLink to next btn
	if pgNum == maxPages {
		pag = append(pag, PaginationLink{Link: "", Name: "→"})
	} else {
		n = strconv.Itoa(pgNum + 1)
		pag = append(pag, PaginationLink{Link: n, Name: "→"})
	}

	return pag
}
