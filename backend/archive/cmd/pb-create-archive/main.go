package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/fsufitch/premium-ball/archive"
	"google.golang.org/protobuf/encoding/protojson"
)

const pokemonTCGSearchURL = "https://api.pokemontcg.io/v2/cards"

func main() {
	maxResults := flag.Int("max", 0, "maximum number of results to retrieve; a number less than 1 means all of them")
	flag.Parse()

	targetZip := flag.Arg(0)
	if targetZip == "" {
		panic(fmt.Errorf("output zipfile not specified"))
	}

	if !strings.HasSuffix(strings.ToLower(targetZip), ".zip") {
		panic(fmt.Errorf("output zipfile must have .zip suffix: %s", targetZip))
	}

	if stat, err := os.Stat(targetZip); err == nil && !stat.Mode().IsRegular() {
		panic(fmt.Errorf("output zip exists but is not a regular file: %s; %w", targetZip, err))
	}

	zipFile, err := os.OpenFile(targetZip, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Errorf("could not open output zip file '%s': %w", targetZip, err))
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	knownCardIDs := map[string]struct{}{}
	for result := range archive.GetAllPokemonTCGCards(pokemonTCGSearchURL) {
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "error fetching cards: %s\n", result.Error)
			continue
		}

		nSuffix := 0
		uniqueID := buildUniqueID(result.Card.Name, result.Card.ID, nSuffix)
		for _, ok := knownCardIDs[uniqueID]; ok; {
			nSuffix++
			uniqueID = buildUniqueID(result.Card.Name, result.Card.ID, nSuffix)
		}
		knownCardIDs[uniqueID] = struct{}{}

		fmt.Printf("%5d. [%s] %s (#%s @ %s)\n", len(knownCardIDs), uniqueID, result.Card.Name, result.Card.Number, result.Card.Set.Name)

		cardData, err := protojson.MarshalOptions{
			Multiline:      true,
			Indent:         "  ",
			UseEnumNumbers: false,
		}.Marshal(archive.ToProto(*result.Card))
		if err != nil {
			panic(fmt.Errorf("could not marshal card: %v", result.Card))
		}

		if w, err := zipWriter.Create(uniqueID + ".json"); err != nil {
			panic(fmt.Errorf("could not add card to archive '%s.json': %w", uniqueID, err))
		} else if _, err := w.Write(cardData); err != nil {
			panic(fmt.Errorf("could not write card to archive '%s.json': %w", uniqueID, err))
		} else if err := zipWriter.Flush(); err != nil {
			panic(fmt.Errorf("could not flush card to archive '%s.json': %w", uniqueID, err))
		}

		if (*maxResults > 0) && len(knownCardIDs) >= *maxResults {
			break
		}
	}

	metaFile, err := zipWriter.Create("metadata.json")
	if err != nil {
		panic(fmt.Errorf("could not create metadata file: %w", err))
	} else if metaData, err := json.MarshalIndent(archiveMetadata{
		RetrievalDate: time.Now().Format("2006/01/02"),
		URL:           pokemonTCGSearchURL,
		Count:         len(knownCardIDs),
	}, "", "  "); err != nil {
		panic(fmt.Errorf("could not marshal metadata json: %w", err))
	} else if _, err := metaFile.Write(metaData); err != nil {
		panic(fmt.Errorf("could not write metadata json to file: %w", err))
	} else if err := zipWriter.Flush(); err != nil {
		panic(fmt.Errorf("could not flush metadata file: %w", err))
	}

	fmt.Fprintf(os.Stderr, "Saved archive of %d cards to: %s\n", len(knownCardIDs), targetZip)
}

func buildUniqueID(cardName string, cardID string, nSuffix int) string {
	idb := &strings.Builder{}
	afterAlphanumeric := false
	for _, ch := range cardName {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			afterAlphanumeric = true
			idb.WriteRune(ch)
		} else if afterAlphanumeric {
			afterAlphanumeric = false
			idb.WriteString("_")
		}
	}
	fmt.Fprintf(idb, "__%s", cardID)
	if nSuffix > 0 {
		fmt.Fprintf(idb, "__%d", nSuffix)
	}
	return idb.String()
}

type archiveMetadata struct {
	RetrievalDate string `json:"retrievalDate"`
	URL           string `json:"url"`
	Count         int    `json:"count"`
}
