package api

import (
	"caxfaxService1/internal/entity"
	"caxfaxService1/pkg/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CatFunFact struct {
	client http.Client
	logger logger.Logger
}

func (c *CatFunFact) GetFunFact() (entity.Fact, error) {
	request, err := http.NewRequest(http.MethodGet, "https://catfact.ninja/fact", nil)
	if err != nil {
		c.logger.Errorf("не удалось создать запрос для получения факта: %s", err.Error())
		return entity.Fact{}, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		c.logger.Errorf("не удалось получить ответ от внешнего API: %s", err.Error())
		return entity.Fact{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.logger.Errorf("не удалось получить ответ от внешнего API: %d", response.StatusCode)
		return entity.Fact{}, err
	}

	rspData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("не удалось получить тело ответа: %s", err.Error())
		return entity.Fact{}, err
	}

	fact := entity.Fact{}
	err = json.Unmarshal(rspData, &fact)
	if err != nil {
		c.logger.Errorf("Не удалось анмаршаллить ответ от API: %s", err.Error())
		return entity.Fact{}, err
	}

	return fact, nil
}

func NewCatFunFact(logger logger.Logger) APIClient {
	cff := &CatFunFact{
		logger: logger,
	}
	cff.client = http.Client{}

	return cff
}
