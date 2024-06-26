package pokeapi

func (c *Client) GetPokemonList(pageURL *string) (PokemonList, error) {
	var resourceResponse PokemonList
	endpoint := baseURL + "/pokemon/"
	if pageURL != nil {
		endpoint = *pageURL
	}

	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	var resourceResponse Pokemon

	endpoint := baseURL + "/pokemon/" + name

	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}

func (c *Client) GetPokemonSpecies(name string) (PokemonSpecies, error) {
	var resourceResponse PokemonSpecies

	endpoint := baseURL + "/pokemon-species/" + name

	err := c.GetJSON(endpoint, &resourceResponse)
	if err != nil {
		return resourceResponse, err
	}
	return resourceResponse, nil
}
