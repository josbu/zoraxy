package acmedns

import (
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"

	"imuslab.com/zoraxy/mod/utils"
)

//go:embed providers.json
var providers []byte //A list of providers generated by acmedns code-generator

type ConfigTemplate struct {
	Name             string `json:"Name"`
	ConfigableFields []struct {
		Title    string `json:"Title"`
		Datatype string `json:"Datatype"`
	} `json:"ConfigableFields"`
	HiddenFields []struct {
		Title    string `json:"Title"`
		Datatype string `json:"Datatype"`
	} `json:"HiddenFields"`
}

// Return a map of string => datatype
func GetProviderConfigStructure(providerName string) (map[string]string, error) {
	//Load the target config template from embedded providers.json
	configTemplateMap := map[string]ConfigTemplate{}
	err := json.Unmarshal(providers, &configTemplateMap)
	if err != nil {
		return map[string]string{}, err
	}

	targetConfigTemplate, ok := configTemplateMap[providerName]
	if !ok {
		return map[string]string{}, errors.New("provider not supported")
	}

	results := map[string]string{}
	for _, field := range targetConfigTemplate.ConfigableFields {
		results[field.Title] = field.Datatype
	}

	return results, nil
}

// HandleServeProvidersJson return the list of supported providers as json
func HandleServeProvidersJson(w http.ResponseWriter, r *http.Request) {
	providerName, _ := utils.GetPara(r, "name")
	if providerName == "" {
		//Send the current list of providers
		configTemplateMap := map[string]ConfigTemplate{}
		err := json.Unmarshal(providers, &configTemplateMap)
		if err != nil {
			utils.SendErrorResponse(w, "failed to load DNS provider")
			return
		}

		//Parse the provider names into an array
		providers := []string{}
		for providerName, _ := range configTemplateMap {
			providers = append(providers, providerName)
		}

		js, _ := json.Marshal(providers)
		utils.SendJSONResponse(w, string(js))
		return
	}
	//Get the config for that provider
	confTemplate, err := GetProviderConfigStructure(providerName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error())
		return
	}

	js, _ := json.Marshal(confTemplate)
	utils.SendJSONResponse(w, string(js))
}