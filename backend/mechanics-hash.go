package premiumball

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/fsufitch/premium-ball/proto"
	"github.com/google/uuid"
)

func CalculateMechanicsHash(card *proto.PokemonTCGCard) string {
	vals := url.Values{}

	vals.Add("name", card.Name)
	vals.Add("supertype", card.Supertype)

	for i, subtype := range sortedStrings(card.Subtypes) {
		vals.Add(fmt.Sprintf("subtype[%d]", i), subtype)
	}

	vals.Add("level", card.Level)
	vals.Add("hp", card.Hp)

	for i, typ := range sortedStrings(card.Types) {
		vals.Add(fmt.Sprintf("type[%d]", i), typ)
	}

	for i, rule := range card.Rules {
		vals.Add(fmt.Sprintf("rule[%d]", i), rule)
	}

	vals.Add("ancientTraitName", card.AncientTrait.Name)
	vals.Add("ancientTraitText", card.AncientTrait.Text)

	for i, ability := range card.Abilities {
		vals.Add(fmt.Sprintf("ability[%d]name", i), ability.Name)
		vals.Add(fmt.Sprintf("ability[%d]text", i), ability.Text)
		vals.Add(fmt.Sprintf("ability[%d]type", i), ability.Type)
	}

	for i, attack := range card.Attacks {
		vals.Add(fmt.Sprintf("attack[%d]name", i), attack.Name)
		for j, cost := range sortedStrings(attack.Cost) {
			vals.Add(fmt.Sprintf("attack[%d]cost[%d]", i, j), cost)
		}
		vals.Add(fmt.Sprintf("attack[%d]damage", i), attack.Damage)
		vals.Add(fmt.Sprintf("attack[%d]text", i), attack.Text)
	}

	for i, weakness := range card.Weaknesses {
		vals.Add(fmt.Sprintf("weakness[%d]", i), weakness.Type+" "+weakness.Value)
	}

	for i, resistance := range card.Resistances {
		vals.Add(fmt.Sprintf("resistance[%d]", i), resistance.Type+" "+resistance.Value)
	}

	for i, rCost := range sortedStrings(card.RetreatCost) {
		vals.Add(fmt.Sprintf("retreatCost[%d]", i), rCost)
	}

	return uuid.NewSHA1(uuid.Nil, []byte(vals.Encode())).String()

	// hashBytes := md5.New().Sum([]byte(vals.Encode()))
	// hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
	// hashBase64 = strings.ToLower(hashBase64)
	// return hashBase64
}

func sortedStrings(s []string) []string {
	s2 := append([]string{}, s...)
	sort.Strings(s2)
	return s2
}
