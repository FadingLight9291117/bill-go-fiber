package lib

import (
	"bill/model"
	"encoding/csv"
	"github.com/samber/lo"
	"os"
)

func ReadCsv(path string, title bool) ([]model.Bill, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	csvReader := csv.NewReader(file)
	dataStr, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if title {
		dataStr = dataStr[1:]
	}
	billList := lo.Map(dataStr, func(item []string, _ int) model.Bill {
		return model.Bill{
			Date:    item[0],
			Money:   item[1],
			Cls:     item[2],
			Label:   item[3],
			Options: item[4],
		}
	})
	return billList, nil
}
