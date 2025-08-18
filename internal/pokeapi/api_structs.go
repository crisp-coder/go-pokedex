package pokeapi

// Utility
type Language struct {
	Id       int
	Name     string
	Official bool
	Iso639   string
	Iso3166  string
	Names    []Name
}

// Common Models
type APIResource struct {
	Url string
}

type APIResourceList struct {
	Count    int
	Next     string
	Previous string
	Results  []APIResource
}

type Description struct {
	Description string
	Language    Language
}

type Effect struct {
	Effect   string
	Language Language
}

type Encounter struct {
	Min_level        int                       `json:"min_level"`
	Max_level        int                       `json:"max_level"`
	Condition_values []EncounterConditionValue `json:"condition_values"`
	Chance           int                       `json:"chance"`
	Method           NamedAPIResource          `json:"method"`
}

type FlavorText struct {
	Flavor_text string
	Language    Language
	Version     Version
}

type GenerationGameIndex struct {
	Game_index int
	Generation Generation
}

type MachineVersionDetail struct {
	//machine      Machine
	VersionGroup VersionGroup
}

type Name struct {
	Name     string
	Language Language
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type NamedAPIResourceList struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

type VerboseEffect struct {
	Effect       string
	Short_effect string
	Language     Language
}

type VersionEncounterDetail struct {
	Max_level         int         `json:"max_level"`
	Version           Version     `json:"version"`
	Max_chance        int         `json:"max_chance"`
	Encounter_details []Encounter `json:"encounter_details"`
}

type VersionGameIndex struct {
	Game_index int
	Version    Version
}

type VersionGroupFlavorText struct {
	Text          string
	Language      Language
	Version_group VersionGroup
}

// ///////////////////////////////
type LocationArea struct {
	Id                     int                   `json:"id"`
	Name                   string                `json:"name"`
	Game_index             int                   `json:"game_index"`
	Encounter_method_rates []EncounterMethodRate `json:"encounter_method_rates"`
	Location               NamedAPIResource      `json:"location"`
	Names                  []Name                `json:"names"`
	Pokemon_encounters     []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	Encounter_method NamedAPIResource         `json:"encounter_method"`
	Version_details  []EncounterVersionDetail `json:"version_details"`
}

type EncounterMethod struct {
	Id    int
	Name  string
	Order int
	Names []Name
}

type EncounterConditionValue struct {
	//Id        int
	//Name      string
	//Condition EncounterCondition
	//Names     []Name
}

type EncounterCondition struct {
	Id     int
	Name   string
	Names  []Name
	Values []EncounterConditionValue
}

type EncounterVersionDetail struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type Machine struct {
	Id int
	//item          Item
	Move          Move
	Version_group VersionGroup
}

type Generation struct {
	Id              int
	Name            string
	Abilities       []Ability
	Names           []Name
	Main_region     Region
	Moves           []Move
	Pokemon_species []PokemonSpecies
	Types           []Type
	Version_groups  []VersionGroup
}

type Ability struct {
}

type Move struct {
}

type Type struct {
}

type Pokedex struct {
	Id              int
	Name            string
	Is_main_series  bool
	Descriptions    []Description
	Names           []Name
	Pokemon_entries []PokemonEntry
	Region          Region
	Version_groups  VersionGroup
}

type Version struct {
	Id            int
	Name          string
	Names         []Name
	Version_group VersionGroup
}

type PokemonEntry struct {
	Entry_number    int
	Pokemon_species PokemonSpecies
}

type PokemonSpecies struct {
}

type VersionGroup struct {
	Id                 int
	Name               string
	Order              int
	Generation         Generation
	Move_learn_methods []MoveLearnMethod
	Pokedexes          []Pokedex
	Regions            []Region
	Versions           []Version
}

type PokemonEncounter struct {
	Pokemon         NamedAPIResource         `json:"pokemon"`
	Version_details []VersionEncounterDetail `json:"version_details"`
}

type Pokemon struct {
	Id              int
	Name            string `json:"name"`
	Base_experience int
	Height          int
	Is_default      bool
	Order           int
	Weight          int
}

type Location struct {
	Id           int                   `json:"id"`
	Name         string                `json:"name"`
	Region       Region                `json:"region"`
	Names        []Name                `json:"names"`
	Game_indices []GenerationGameIndex `json:"game_indices"`
	Areas        []LocationArea        `json:"areas"`
}

type Region struct {
	Id        int
	Locations []Location
	Name      string
	Names     []Name
	//Main_generation Generation
	Pokedexes      []Pokedex
	Version_groups []VersionGroup
}

type MoveLearnMethod struct {
	Id             int
	Name           string
	Descriptions   []Description
	Names          []Name
	Version_groups []VersionGroup
}
