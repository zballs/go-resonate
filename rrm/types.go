package rrm

import (
	"encoding/json"
)

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"
)

// json-ld types for the coalaip rrm
// specs: github.com/COALAIP/specs/data-structure/README.md
// TODO: ipld

// Geo coordinates

func GeoData(lat, long string) []byte {
	data, _ := json.Marshal(struct {
		Context   string `json:"@context"`
		Type      string `json:"@type"`
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}{
		Context:   SCHEMA,
		Type:      "GeoCoordinates",
		Latitude:  lat,
		Longitude: long,
	})
	return data
}

func GeoMap(lat, long string) map[string]interface{} {
	data := GeoData(lat, long)
	mp, _ := MapJSON(data)
	return mp
}

// Localizable place

func PlaceData(lat, long, name string) []byte {
	data, _ := json.Marshal(struct {
		Context string                 `json:"@context"`
		Type    string                 `json:"@type"`
		Geo     map[string]interface{} `json:"geo"`
		Name    string                 `json:"name"`
	}{
		Context: SCHEMA,
		Type:    "Place",
		Geo:     GeoMap(lat, long),
		Name:    name,
	})
	return data
}

// Person

func PersonData(id, givenName, familyName, birthDate, deathDate string) []byte {
	data, _ := json.Marshal(struct {
		Context   string `json: "@context"`
		Type      string `json:"@type"`
		Id        string `json:"@id"`
		Email     string `json:"email"`
		BirthDate string `json:"birthDate"`
	}{
		Context:    SCHEMA,
		Type:       "Person",
		Id:         id,
		GivenName:  givenName,
		FamilyName: familyName,
		BirthDate:  birthDate,
		DeathDate:  deathDate,
	})
	return data
}

// Organization

func OrganizationData(id, name string, founder string, members []string) []byte {
	data, _ := json.Marshal(struct {
		Context string   `json:"@context"`
		Type    string   `json:"@type"`
		Id      string   `json:"@id"`
		Name    string   `json:"name"`
		Founder string   `json:"founder"`
		Member  []string `json:"member"`
	}{
		Context: SCHEMA,
		Type:    "Organization",
		Id:      id,
		Name:    name,
		Founder: founder,
		Member:  members,
	})
	return data
}

//---------------------------------------------------

// Creation

func WorkData(id, name, author string) []byte {
	data, _ := json.Marshal(struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Id      string `json:"@id"`
		Name    string `json:"name"`
		Author  string `json:"author"`
	}{
		Context: COALAIP,
		Type:    "Work",
		Id:      id,
		Name:    name,
		Author:  author,
	})
	return data
}

// Digital manifestation

func DigitalManifestationData(id string, name string, example string, isManifestation bool, project string, datePublished string, locationCreated, url string) []byte {
	data, _ := json.Marshal(struct {
		Context         string `json:"@context"`
		Type            string `json:"@type"`
		Id              string `json:"@id"`
		Name            string `json:"name"`
		ExampleOfWork   string `json:"exampleOfWork"`
		IsManifestation bool   `json:"isManifestation"`
		IsPartOf        string `json:"isPartOf"`
		DatePublished   string `json:"datePublished"`
		LocationCreated string `json:"locationCreated"`
		Url             string `json:"url"`
	}{
		Context:         COALAIP,
		Type:            "Manifestation",
		Id:              id,
		Name:            name,
		ExampleOfWork:   example,
		IsManifestation: isManifestation,
		IsPartOf:        project,
		DatePublished:   datePublished,
		LocationCreated: locationCreated,
		Url:             url,
	})
	return data
}

// Digital Fingerprint

func digitalFingerprint(id, creativeWork string, fingerprint string) []byte {
	data, _ := json.Marshal(struct {
		Context       string `json:"@context"`
		Type          string `json:"@type"`
		Id            string `json:"@id"`
		FingerprintOf string `json:"fingerprintOf"`
		Fingerprint   string `json:"fingerprint"`
	}{
		Context:       COALAIP,
		Type:          "DigitalFingerprint",
		Id:            id,
		FingerprintOf: creativeWork,
		Fingerprint:   fingerprint,
	})
	return data
}

// Right

func RightData(id string, usages []string, territory string, rightContext []string, exclusive bool, numberOfUses, percentageShares int, validFrom, validTo string, creativeWork, license string) []byte {
	data, _ := json.Marshal(struct {
		Context          string   `json:"@context"`
		Type             string   `json:"@type"`
		Id               string   `json:"@id"`
		Usages           []string `json:"usages"`
		Territory        string   `json:"territory"`
		RightContext     []string `json:"rightContext"`
		Exclusive        bool     `json:"exclusive"`
		NumberOfUses     int      `json:"numberOfUses"`
		PercentageShares int      `json:"share"`
		ValidFrom        string   `json:"validFrom"`
		ValidTo          string   `json:"validTo"`
		Creation         string   `json:"creation"`
		License          string   `json:"license"`
	}{
		Context:          COALAIP,
		Type:             "Right",
		Id:               id,
		Usages:           usages,
		Territory:        territory,
		RightContext:     rightContext,
		Exclusive:        exclusive,
		NumberOfUses:     numberOfUses,
		PercentageShares: percentageShares,
		ValidFrom:        validFrom,
		ValidTo:          validTo,
		Creation:         creativeWork,
		License:          license,
	})
	return data
}

// Rights assignment

func RightsAssignmentData(id, creativeWork string) []byte {
	data, _ := json.Marshal(struct {
		Context          string `json:"@context"`
		Type             string `json:"@type"`
		Id               string `json:"@id"`
		TransferContract string `json:"transferContract"`
	}{
		Context:          COALAIP,
		Type:             "RightsTransferAction",
		Id:               id,
		TransferContract: creativeWork,
	})
	return data
}

// Rights assertion

func RightsAssertionData(id, asserter string, assertionTruth bool, assertionSubject string, _error, validFrom, validThrough string) []byte {
	data, _ := json.Marshal(struct {
		Context          string `json:"@context"`
		Type             string `json:"@type"`
		Id               string `json:"@id"`
		Asserter         string `json:"asserter"`
		AssertionTruth   bool   `json:"assertionTruth"`
		AssertionSubject string `json:"assertionSubject"`
		Error            string `json:"error"`
		ValidFrom        string `json:"validFrom"`
		ValidThrough     string `json:"validThrough"`
	}{
		Context:          SCHEMA,
		Type:             "ReviewAction",
		Id:               id,
		Asserter:         asserter,
		AssertionTruth:   assertionTruth,
		AssertionSubject: assertionSubject,
		Error:            _error,
		ValidFrom:        validFrom,
		ValidThrough:     validThrough,
	})
	return data
}
