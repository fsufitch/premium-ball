package archive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/fsufitch/premium-ball/proto"
)

type tcgCardSearchResult struct {
	Card  *PTCGCard
	Error error
}

// GetAllPokemonTCGCards returns asynchronous channels that feed through all cards in the database (as encoding/json unmarshaled values), or errors encountered while querying
func GetAllPokemonTCGCards(searchURL string) (chOutput chan tcgCardSearchResult) {
	chOutput = make(chan tcgCardSearchResult)
	go func() {
		defer close(chOutput)
		numPages := -1
		page := 1
		baseURL, err := url.Parse(searchURL)
		if err != nil {
			chOutput <- tcgCardSearchResult{nil, fmt.Errorf("bad base search URL (%s): %w", searchURL, err)}
			return
		}
		for numPages < 0 || page <= numPages {
			var url url.URL = *baseURL

			qvals := url.Query()
			qvals.Set("page", fmt.Sprintf("%d", page))
			url.RawQuery = qvals.Encode()
			fmt.Println(url.String())

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
			page++
		}
	}()
	return
}

func ToProto(card PTCGCard) *proto.PokemonTCGCard {
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
			Id:           card.Set.ID,
			Name:         card.Set.Name,
			Series:       card.Set.Series,
			PrintedTotal: card.Set.PrintedTotal,
			Total:        card.Set.Total,
			Legalities: &proto.Legalities{
				Standard:  card.Set.Legalities.Standard,
				Expanded:  card.Set.Legalities.Expanded,
				Unlimited: card.Set.Legalities.Unlimited,
			},
			PtcgoCode:   card.Set.PTCGOCode,
			ReleaseDate: card.Set.ReleaseDate,
			UpdatedAt:   card.Set.UpdatedAt,
			Images: &proto.Images{
				Small: card.Set.Images.Small,
				Large: card.Set.Images.Large,
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
			Url:       card.TCGPlayer.URL,
			Updatedat: card.TCGPlayer.UpdatedAt,
			PricesUSD: map[string]*proto.TCGPlayerPricesUSD{}, // fill below
		},
		CardMarket: &proto.CardMarketDetails{
			Url:       card.CardMarket.URL,
			UpdatedAt: card.CardMarket.UpdatedAt,
			PricesEUR: &proto.CardMarketPricesEUR{
				AverageSellPrice: card.CardMarket.Prices.AverageSellPrice,
				LowPrice:         card.CardMarket.Prices.LowPrice,
				TrendPrice:       card.CardMarket.Prices.TrendPrice,
				GermanProLow:     card.CardMarket.Prices.GermanProLow,
				SuggestedPrice:   card.CardMarket.Prices.SuggestedPrice,
				ReverseHoloSell:  card.CardMarket.Prices.ReverseHoloSell,
				ReverseHoloLow:   card.CardMarket.Prices.ReverseHoloLow,
				ReverseHoloTrend: card.CardMarket.Prices.ReverseHoloTrend,
				LowPriceExPlus:   card.CardMarket.Prices.LowPriceExPlus,
				Avg1:             card.CardMarket.Prices.Avg1,
				Avg7:             card.CardMarket.Prices.Avg7,
				Avg30:            card.CardMarket.Prices.Avg30,
				ReverseHoloAvg1:  card.CardMarket.Prices.Avg1,
				ReverseHoloAvg7:  card.CardMarket.Prices.ReverseHoloAvg7,
				ReverseHoloAvg30: card.CardMarket.Prices.Avg30,
			},
		},
	}

	for _, ability := range card.Abilities {
		pbCardData.Abilities = append(pbCardData.Abilities, &proto.Ability{
			Name: ability.Name,
			Text: ability.Text,
			Type: ability.Type,
		})
	}

	for _, attack := range card.Attacks {
		pbCardData.Attacks = append(pbCardData.Attacks, &proto.Attack{
			Cost:                append([]string{}, attack.Cost...),
			Name:                attack.Name,
			Text:                attack.Text,
			Damage:              attack.Damage,
			ConvertedEnergyCost: int64(attack.ConvertedEnergyCost),
		})
	}

	for _, weakness := range card.Weaknesses {
		pbCardData.Weaknesses = append(pbCardData.Weaknesses, &proto.Weakness{
			Type:  weakness.Type,
			Value: weakness.Value,
		})
	}

	for _, resistance := range card.Resistances {
		pbCardData.Resistances = append(pbCardData.Resistances, &proto.Resistance{
			Type:  resistance.Type,
			Value: resistance.Value,
		})
	}

	for _, nationalPokedexNumber := range card.NationalPokedexNumbers {
		pbCardData.NationalPokedexNumbers = append(pbCardData.NationalPokedexNumbers, int64(nationalPokedexNumber))
	}

	for k, v := range card.TCGPlayer.Prices {
		pbCardData.TcgPlayer.PricesUSD[k] = &proto.TCGPlayerPricesUSD{
			Low:       v.Low,
			Mid:       v.Mid,
			High:      v.High,
			Market:    v.Market,
			DirectLow: v.DirectLow,
		}
	}

	return pbCardData
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
	RetreatCost          []string `json:"retreatCost"`
	ConvertedRetreatCost int      `json:"convertedRetreatCost"`
	Set                  struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Series       string `json:"series"`
		PrintedTotal int64  `json:"printedTotal"`
		Total        int64  `json:"total"`
		Legalities   struct {
			Standard  string `json:"standard"`
			Expanded  string `json:"expanded"`
			Unlimited string `json:"unlimited"`
		} `json:"legalities"`
		PTCGOCode   string `json:"ptcgoCode"`
		ReleaseDate string `json:"releaseDate"`
		UpdatedAt   string `json:"updatedAt"`
		Images      struct {
			Small string `json:"small"`
			Large string `json:"large"`
		} `json:"images"`
	} `json:"set"`
	Number                 string `json:"number"`
	Artist                 string `json:"artist"`
	Rarity                 string `json:"rarity"`
	FlavorText             string `json:"flavorText"`
	NationalPokedexNumbers []int  `json:"nationalPokedexNumbers"`
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

	TCGPlayer struct {
		URL       string `json:"url"`
		UpdatedAt string `json:"updatedAt"`
		Prices    map[string]struct {
			Low       float32 `json:"low"`
			Mid       float32 `json:"mid"`
			High      float32 `json:"high"`
			Market    float32 `json:"market"`
			DirectLow float32 `json:"directLow"`
		} `json:"prices"`
	} `json:"tcgplayer"`
	CardMarket struct {
		URL       string `json:"url"`
		UpdatedAt string `json:"updatedAt"`
		Prices    struct {
			AverageSellPrice float32 `json:"averageSellPrice"`
			LowPrice         float32 `json:"lowPrice"`
			TrendPrice       float32 `json:"trendPrice"`
			GermanProLow     float32 `json:"germanProLow"`
			SuggestedPrice   float32 `json:"suggestedPrice"`
			ReverseHoloSell  float32 `json:"reverseHoloSell"`
			ReverseHoloLow   float32 `json:"reverseHoloLow"`
			ReverseHoloTrend float32 `json:"reverseHoloTrend"`
			LowPriceExPlus   float32 `json:"lowPriceExPlus"`
			Avg1             float32 `json:"avg1"`
			Avg7             float32 `json:"avg7"`
			Avg30            float32 `json:"avg30"`
			ReverseHoloAvg1  float32 `json:"reverseHoloAvg1"`
			ReverseHoloAvg7  float32 `json:"reverseHoloAvg7"`
			ReverseHoloAvg30 float32 `json:"reverseHoloAvg30"`
		} `json:"prices"`
	} `json:"cardmarket"`
}
