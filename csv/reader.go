package csv

import (
	"encoding/csv"
	"io"
)

func ReadSites(file io.Reader) (map[string]string, error) {
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	sites := map[string]string{}

	for i := 1; i < len(rows); i++ {
		sites[rows[i][0]] = rows[i][1]
	}

	return sites, err
}
