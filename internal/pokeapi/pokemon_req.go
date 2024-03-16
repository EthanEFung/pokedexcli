package pokeapi

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	var resourceResponse Pokemon

	endpoint := baseURL + "/pokemon/" + name

	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}
