package helpers

import "strconv"

type PaginationLink struct {
	Link string
	Name string
}

func GeneratePagination(pgNum int, maxPages int) []PaginationLink {
	pag := make([]PaginationLink, 0)
	var n string
	delta := 2

	// Add pagLink to prev btn
	if pgNum == 1 {
		pag = append(pag, PaginationLink{Link: "", Name: "←"})
	} else {
		n = strconv.Itoa(pgNum - 1)
		pag = append(pag, PaginationLink{Link: n, Name: "←"})
	}

	// Simple fill slice for 6 pages
	if maxPages <= 6 {
		for i := 1; i <= maxPages; i++ {
			n = strconv.Itoa(i)
			pag = append(pag, PaginationLink{Link: n, Name: n})
		}
	} else {
		// Fill slice for a lot of pages
		switch pgNum {
		case 1:
			for i := pgNum; i <= pgNum+delta*2; i++ {
				n = strconv.Itoa(i)
				pag = append(pag, PaginationLink{Link: n, Name: n})
			}

			pag = append(pag, PaginationLink{Link: "", Name: "..."})
		case 2:
			for i := pgNum - 1; i <= pgNum+delta; i++ {
				n = strconv.Itoa(i)
				pag = append(pag, PaginationLink{Link: n, Name: n})
			}

			pag = append(pag, PaginationLink{Link: "", Name: "..."})
		case maxPages - 1:
			pag = append(pag, PaginationLink{Link: "", Name: "..."})

			for i := pgNum - delta; i <= maxPages; i++ {
				n = strconv.Itoa(i)
				pag = append(pag, PaginationLink{Link: n, Name: n})
			}
		case maxPages:
			pag = append(pag, PaginationLink{Link: "", Name: "..."})

			for i := pgNum - 2*delta; i <= maxPages; i++ {
				n = strconv.Itoa(i)
				pag = append(pag, PaginationLink{Link: n, Name: n})
			}
		default:
			if pgNum-delta != 1 {
				pag = append(pag, PaginationLink{Link: "", Name: "..."})
			}

			for i := pgNum - delta; i <= pgNum+delta; i++ {
				n = strconv.Itoa(i)
				pag = append(pag, PaginationLink{Link: n, Name: n})
			}

			if pgNum+delta != maxPages {
				pag = append(pag, PaginationLink{Link: "", Name: "..."})
			}
		}
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
