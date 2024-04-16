package namefinder

type BasicPokemonInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Form string `json:"form,omitempty"`
}

type BasicPokemonInfoEntries [][]BasicPokemonInfo
