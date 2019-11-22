package service

import (
	"encoding/json"
	"time"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

type UsageService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *UsageService) GetUsage(startDate, endDate float64) {
	params := &operations.GetUsageParams{
		StartDate: startDate,
		EndDate:   endDate,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetUsage(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	} else {
		jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
		if err != nil {
			log.Fatalln("err: ", err)
		}
		log.Infoln(string(jsonPayload))
	}
}
