package pokeapi

func (c *Client) GetLocations(pageURL *string) (LocationResponse, error) {
	var resourceResponse LocationResponse
	endpoint := baseURL + "/location"
	if pageURL != nil {
		endpoint = *pageURL
	}
	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}

func (c *Client) GetLocation(name string) (Location, error) {
	var resourceResponse Location
	endpoint := baseURL + "/location/" + name
	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}
