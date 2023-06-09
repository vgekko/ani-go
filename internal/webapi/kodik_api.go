package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vgekko/ani-go/internal/entity"
	"github.com/vgekko/ani-go/pkg/apperror"
)

type KodikWebAPI struct {
	token  string
	client http.Client
}

func NewKodikWebAPI() *KodikWebAPI {
	token := os.Getenv("KODIK_TOKEN")
	client := http.Client{Timeout: time.Second * 3}

	return &KodikWebAPI{token: token, client: client}
}

func (k *KodikWebAPI) SearchTitles(option, value string) (entity.KodikAPI, error) {
	var kodikResponse entity.KodikAPI

	url := fmt.Sprintf("https://kodikapi.com/search?token=%s&%s=%s", k.token, option, value)

	resp, err := k.client.Get(url)
	if err != nil {
		return entity.KodikAPI{}, fmt.Errorf("webapi.ResultsByKinopoiskID: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&kodikResponse); err != nil {
		return entity.KodikAPI{}, err
	}

	if kodikResponse.Total == 0 {
		return entity.KodikAPI{}, apperror.ErrTitleNotFound
	}

	return kodikResponse, nil
}
