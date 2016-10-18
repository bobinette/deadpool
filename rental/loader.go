package rental

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func Load(p Parameters) (StateRewards, bool, error) {
	filename := fmt.Sprintf("%d_%d_%d_%d_%d_%d.json", p.CustomerAt1, p.CustomerAt2, p.ReturnAt1, p.ReturnAt2, p.MaxCars, p.MaxMoves)
	filepath := path.Join(".", "rental", "data", filename)
	file, err := os.Open(filepath)
	if err == os.ErrNotExist {
		return nil, false, nil
	} else if err != nil {
		return nil, true, err
	}
	defer file.Close()

	r := make(StateRewards)
	err = json.NewDecoder(file).Decode(&r)
	if err != nil {
		return nil, true, err
	}

	return r, true, nil
}

func Save(p Parameters, r StateRewards) error {
	filename := fmt.Sprintf("%d_%d_%d_%d_%d_%d.json", p.CustomerAt1, p.CustomerAt2, p.ReturnAt1, p.ReturnAt2, p.MaxCars, p.MaxMoves)
	filepath := path.Join(".", "rental", "data", filename)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(r)
}
