package coala

import . "github.com/zballs/go_resonate/util"

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"
)

// json-ld types for the coalaip rrm
// specs: github.com/COALAIP/specs/data-structure/README.md
// Note: @id should be set after model is created
// TODO: ipld implementation

// Geo coordinates

func NewGeo(lat, long string) map[string]interface{} {
	json := MarshalJSON(struct {
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
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Localizable place

func NewPlace(lat, long, name string) map[string]interface{} {
	json := MarshalJSON(struct {
		Context string                 `json:"@context"`
		Type    string                 `json:"@type"`
		Geo     map[string]interface{} `json:"geo"`
		Name    string                 `json:"name"`
	}{
		Context: SCHEMA,
		Type:    "Place",
		Geo:     NewGeo(lat, long),
		Name:    name,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Person

func NewPerson(givenName, familyName, birthDate string) map[string]interface{} {
	json := MarshalJSON(struct {
		Context    string `json: "@context"`
		Type       string `json:"@type"`
		Id         string `json:"@id"`
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
		BirthDate  string `json:"birthDate"`
	}{
		Context:    SCHEMA,
		Type:       "Person",
		Id:         "",
		GivenName:  givenName,
		FamilyName: familyName,
		BirthDate:  birthDate,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Organization

func NewOrganization(name string, founder string, members []string) map[string]interface{} {
	json := MarshalJSON(struct {
		Context string   `json:"@context"`
		Type    string   `json:"@type"`
		Id      string   `json:"@id"`
		Name    string   `json:"name"`
		Founder string   `json:"founder"`
		Member  []string `json:"member"`
	}{
		Context: SCHEMA,
		Type:    "Organization",
		Name:    name,
		Founder: founder,
		Member:  members,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Work (e.g. composition)

func NewWork(name, author string) map[string]interface{} {
	json := MarshalJSON(struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Id      string `json:"@id"`
		Name    string `json:"name"`
		Author  string `json:"author"`
	}{
		Context: COALAIP,
		Type:    "Work",
		Name:    name,
		Author:  author,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Digital manifestation (e.g. audio file)

func NewDigitalManifestation(name string, example string, isManifestation bool, project string, datePublished string, locationCreated, url string) map[string]interface{} {
	json := MarshalJSON(struct {
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
		Name:            name,
		ExampleOfWork:   example,
		IsManifestation: isManifestation,
		IsPartOf:        project,
		DatePublished:   datePublished,
		LocationCreated: locationCreated,
		Url:             url,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Digital Fingerprint

func NewDigitalFingerprint(creativeWork string, fingerprint string) map[string]interface{} {
	json := MarshalJSON(struct {
		Context       string `json:"@context"`
		Type          string `json:"@type"`
		Id            string `json:"@id"`
		FingerprintOf string `json:"fingerprintOf"`
		Fingerprint   string `json:"fingerprint"`
	}{
		Context:       COALAIP,
		Type:          "DigitalFingerprint",
		FingerprintOf: creativeWork,
		Fingerprint:   fingerprint,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Right

func NewRight(usages []string, territory string, rightContext []string, exclusive bool, numberOfUses, percentageShares int, validFrom, validTo string, creativeWork, license string) map[string]interface{} {
	json := MarshalJSON(struct {
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
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Rights assignment

func NewRightsAssignment(id, creativeWork string) map[string]interface{} {
	json := MarshalJSON(struct {
		Context          string `json:"@context"`
		Type             string `json:"@type"`
		Id               string `json:"@id"`
		TransferContract string `json:"transferContract"`
	}{
		Context:          COALAIP,
		Type:             "RightsTransferAction",
		TransferContract: creativeWork,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}

// Rights assertion

func NewRightsAssertion(id, asserter string, assertionTruth bool, assertionSubject string, _error, validFrom, validThrough string) map[string]interface{} {
	json := MarshalJSON(struct {
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
		Asserter:         asserter,
		AssertionTruth:   assertionTruth,
		AssertionSubject: assertionSubject,
		Error:            _error,
		ValidFrom:        validFrom,
		ValidThrough:     validThrough,
	})
	data, err := CompactJSON(json)
	Check(err)
	return data
}
