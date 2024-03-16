package pokeapi

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	var resourceResponse LocationAreaResponse

	endpoint := baseURL + "/location-area"
	if pageURL != nil {
		endpoint = *pageURL
	}

	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}

func (c *Client) GetLocationArea(areaName string) (LocationArea, error) {

	var resourceResponse LocationArea

	endpoint := baseURL + "/location-area/" + areaName

	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}
