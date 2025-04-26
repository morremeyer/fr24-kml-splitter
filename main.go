package main

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type rawXML struct {
	Inner []byte `xml:",innerxml"`
}

type kml struct {
	Document document `xml:"Document"`
	Xmlns    string   `xml:"xmlns,attr"`
}

type document struct {
	Name        rawXML   `xml:"name"`
	Description rawXML   `xml:"description"`
	Folders     []folder `xml:"Folder"`
}

type folder struct {
	Name       string      `xml:"name"`
	Placemarks []placemark `xml:"Placemark"`
}

type placemark struct {
	Name        rawXML `xml:"name"`
	Description rawXML `xml:"description"`
	Timestamp   struct {
		When time.Time `xml:"when"`
	} `xml:"TimeStamp"`
	Style rawXML `xml:"Style"`
	Point rawXML `xml:"Point"`
}

func main() {
	argLength := len(os.Args[1:])
	if argLength == 0 {
		log.Fatal("You need to specify the input file")
	}

	input, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	var kmlFile kml
	if err := xml.Unmarshal(data, &kmlFile); err != nil {
		log.Fatal(err)
	}

	// Remove the trail
	trailIdx := slices.IndexFunc(kmlFile.Document.Folders, func(f folder) bool { return f.Name == "Trail" })
	if trailIdx != -1 {
		kmlFile.Document.Folders = slices.Delete(kmlFile.Document.Folders, trailIdx, trailIdx+1)
	}

	// If a timestamp to split is specified, we handle it here
	if argLength == 2 {
		splitTime, err := time.Parse("2006-01-02T15:04:05", os.Args[2])
		if err != nil {
			log.Fatalf("Split time could not be parsed: %s", err)
		}

		timeIdx := slices.IndexFunc(kmlFile.Document.Folders[0].Placemarks, func(p placemark) bool { return p.Timestamp.When.Equal(splitTime) })

		orginalPlacemarks := kmlFile.Document.Folders[0].Placemarks

		// We need a deep copy of the kml
		var folders []folder
		second := kml{
			Xmlns: kmlFile.Xmlns,
			Document: document{
				Name:        kmlFile.Document.Name,
				Description: kmlFile.Document.Description,
				Folders:     folders,
			},
		}

		kmlFile.Document.Folders[0].Placemarks = orginalPlacemarks[:timeIdx]
		second.Document.Folders = append(second.Document.Folders, folder{
			Placemarks: orginalPlacemarks[timeIdx:],
		})

		// Save first
		finalFirst, err := xml.MarshalIndent(kmlFile, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile("route-only-first.kml", finalFirst, 0750); err != nil {
			log.Fatal(err)
		}

		// Save second
		finalSecond, err := xml.MarshalIndent(second, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile("route-only-second.kml", finalSecond, 0750); err != nil {
			log.Fatal(err)
		}

	} else {
		final, err := xml.MarshalIndent(kmlFile, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile("route-only.kml", final, 0750); err != nil {
			log.Fatal(err)
		}
	}
}
