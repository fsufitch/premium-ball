package archive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	premiumball "github.com/fsufitch/premium-ball"
	"github.com/fsufitch/premium-ball/proto"
)

const pokemonTCGSearchURL = "https://api.pokemontcg.io/v2/cards"

type tcgCardSearchResult struct {
	Card  *PTCGCard
	Error error
}

// GetAllPokemonTCGCards returns asynchronous channels that feed through all cards in the database (as encoding/json unmarshaled values), or errors encountered while querying
func GetAllPokemonTCGCards() (chOutput chan tcgCardSearchResult) {
	chOutput = make(chan tcgCardSearchResult)
	go func() {
		defer close(chOutput)
		numPages := -1
		page := 1
		baseURL, _ := url.Parse(pokemonTCGSearchURL)
		for numPages < 0 || page <= numPages {
			var url url.URL = *baseURL
			url.Query().Add("page", fmt.Sprintf("%d", page))
			respBody := PTCGCardSearchResponse{}

			if resp, err := http.Get(url.String()); err != nil {
				chOutput <- tcgCardSearchResult{nil, fmt.Errorf("query failed (url=%s): %w", url.String(), err)}
				return
			} else if body, err := io.ReadAll(resp.Body); err != nil {
				chOutput <- tcgCardSearchResult{nil, fmt.Errorf("failed reading response (url=%s): %w", url.String(), err)}
				return
			} else if err := json.Unmarshal(body, &respBody); err != nil {
				chOutput <- tcgCardSearchResult{nil, fmt.Errorf("failed unmarsharing query body (url=%s): %w", url.String(), err)}
				return
			}

			if numPages < 0 {
				numPages = (respBody.TotalCount / respBody.PageSize) + 1
			}

			for i := range respBody.Data {
				chOutput <- tcgCardSearchResult{&respBody.Data[i], nil}
			}
		}
	}()
	return
}

func ToPremiumBallCard(card PTCGCard) *proto.PremiumBallCard {
	pbCardData := &proto.PokemonTCGCard{
		Id:          card.ID,
		Name:        card.Name,
		Supertype:   card.Supertype,
		Subtypes:    append([]string{}, card.Subtypes...),
		Level:       card.Level,
		Hp:          card.HP,
		Types:       append([]string{}, card.Types...),
		EvolvesFrom: "",
		EvolvesTo:   "",
		Rules:       append([]string{}, card.Rules...),
		AncientTrait: &proto.AncientTrait{
			Name: card.AncientTrait.Name,
			Text: card.AncientTrait.Text,
		},
		Abilities:            []*proto.Ability{},    // fill below
		Attacks:              []*proto.Attack{},     // fill below
		Weaknesses:           []*proto.Weakness{},   // fill below
		Resistances:          []*proto.Resistance{}, // fill below
		RetreatCost:          append([]string{}, card.RetreatCost...),
		ConvertedRetreatCost: int64(card.ConvertedRetreatCost),
		Set: &proto.PokemonTCGSet{
			// XXX: TODO
			Id:           "",
			Name:         "",
			Series:       "",
			PrintedTotal: 0,
			Total:        0,
			Legalities: &proto.Legalities{
				Standard:  "",
				Expanded:  "",
				Unlimited: "",
			},
			PtcgoCode:   "",
			ReleaseDate: "",
			UpdatedAt:   "",
			Images: &proto.Images{
				Small: "",
				Large: "",
			},
		},
		Number:                 card.Number,
		Artist:                 card.Artist,
		Rarity:                 card.Rarity,
		FlavorText:             card.FlavorText,
		NationalPokedexNumbers: []int64{}, // fill below
		Legalities: &proto.Legalities{
			Standard:  card.Legalities.Standard,
			Expanded:  card.Legalities.Expanded,
			Unlimited: card.Legalities.Unlimited,
		},
		RegulationMark: card.RegulationMark,
		Images: &proto.Images{
			Small: card.Images.Small,
			Large: card.Images.Large,
		},
		TcgPlayer: &proto.TCGPlayerDetails{
			// XXX: TODO
			Url:       "",
			Updatedat: "",
			Prices: &proto.TCGPlayerPricesUSD{
				Low:       0.0,
				Mid:       0.0,
				High:      0.0,
				Market:    0.0,
				DirectLow: 0.0,
			},
		},
		CardMarket: &proto.CardMarketPricesUSD{
			// XXX: TODO
			AverageSellPrice: 0.0,
			LowPrice:         0.0,
			TrendPrice:       0.0,
			GermanProLow:     0.0,
			SuggestedPrice:   0.0,
			ReverseHoloSell:  0.0,
			ReverseHoloLow:   0.0,
			ReverseHoloTrend: 0.0,
			LowPriceExPlus:   0.0,
			Avg1:             0.0,
			Avg7:             0.0,
			Avg30:            0.0,
			ReverseHoloAvg1:  0.0,
			ReverseHoloAvg7:  0.0,
			ReverseHoloAvg30: 0.0,
		},
	}

	pbc := &proto.PremiumBallCard{
		Id:            card.ID,
		CardData:      pbCardData,
		MechanicsHash: premiumball.CalculateMechanicsHash(pbCardData),
	}
	return pbc
}

type PTCGCardSearchResponse struct {
	Data       []PTCGCard `json:"data"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	Count      int        `json:"count"`
	TotalCount int        `json:"totalCount"`
}

type PTCGCard struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Supertype    string   `json:"supertype"`
	Subtypes     []string `json:"subtypes"`
	Level        string   `json:"level"`
	HP           string   `json:"hp"`
	Types        []string `json:"types"`
	EvolvesFrom  string   `json:"evolvesFrom"`
	EvolvesTo    []string `json:"evolvesTo"`
	Rules        []string `json:"rules"`
	AncientTrait struct {
		Name string `json:"name"`
		Text string `json:"text"`
	} `json:"ancientTrait"`
	Abilities []struct {
		Name string `json:"name"`
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"abilities"`
	Attacks []struct {
		Cost                []string `json:"cost"`
		Name                string   `json:"name"`
		Text                string   `json:"text"`
		Damage              string   `json:"damage"`
		ConvertedEnergyCost int      `json:"convertedEnergyCost"`
	} `json:"attacks"`
	Weaknesses []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"weaknesses"`
	Resistances []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"resistances"`
	RetreatCost            []string    `json:"retreatCost"`
	ConvertedRetreatCost   int         `json:"convertedRetreatCost"`
	Set                    interface{} `json:"set"`
	Number                 string      `json:"number"`
	Artist                 string      `json:"artist"`
	Rarity                 string      `json:"rarity"`
	FlavorText             string      `json:"flavorText"`
	NationalPokedexNumbers []int       `json:"nationalPokedexNumbers"`
	Legalities             struct {
		Standard  string `json:"standard"`
		Expanded  string `json:"expanded"`
		Unlimited string `json:"unlimited"`
	} `json:"legalities"`
	RegulationMark string `json:"regulationMark"`
	Images         struct {
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"images"`

	TCGPlayer  interface{} `json:"tcgplayer"`
	CardMarket interface{} `json:"cardmarket"`
}
