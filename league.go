package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewLeague(rdr io.Reader) (League, error) {
	var league League

	err := json.NewDecoder(rdr).Decode(&league)

	if err != nil {
		err = fmt.Errorf("unable to parse league, %v", err)
	}

	return league, err
}
