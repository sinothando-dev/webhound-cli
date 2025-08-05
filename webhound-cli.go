package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type DataAttributeDesciption struct {
	Value      string   `json:"value"`
	SourceUrls []string `json:"source_urls"`
}

type DataAttributeImageUrl struct {
	Value      string   `json:"value"`
	SourceUrls []string `json:"source_urls"`
}

type DataAttributeTags struct {
	Value      []string `json:"value"`
	SourceUrls []string `json:"source_urls"`
}

type DataAttributes struct {
	ImageUrl    DataAttributeImageUrl   `json:"image_url"`
	Description DataAttributeDesciption `json:"description"`
	Tags        DataAttributeTags       `json:"tags"`
}

type Data struct {
	Id           string         `json:"id"`
	CreatedAt    string         `json:"created_at"`
	Is_valid     bool           `json:"is_valid"`
	Is_duplicate bool           `json:"is_duplicate"`
	Attributes   DataAttributes `json:"attributes"`
}

type MetaDataAttributes struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Is_primary  bool   `json:"is_primary"`
}

type MetaDataSchema struct {
	Attributes []MetaDataAttributes `json:"attributes"`
}

type MetaDataEntity struct {
	Name        string   `json:"name"`
	Criteria    []string `json:"criteria"`
	Description string   `json:"description"`
}

type MetaData struct {
	Entity      MetaDataEntity `json:"entity"`
	Schema      MetaDataSchema `json:"schema"`
	Exported_At string         `json:"exported_at"`
	TotalItems  int            `json:"total_items"`
}

type WebhoundJson struct {
	MetaData MetaData `json:"metadata"`
	Data     []Data   `json:"data"`
}

func main() {
	json_file_path := flag.String("path", "", "--pathToYourJsonFile")
	folder_name := flag.String("name", "noname", "--yourNewFolderName")
	flag.Parse()
	json_file, err := os.Open(*json_file_path)

	if err != nil {
		log.Fatal("Problem with finding the json file")
	}

	json_byte_value, _ := io.ReadAll(json_file)

	var webhoundJson WebhoundJson

	json.Unmarshal(json_byte_value, &webhoundJson)

	if err != nil {
		log.Fatal("Error when trying to parse the json file")
	}

	if *folder_name == "noname" {
		folder_name = &webhoundJson.MetaData.Entity.Name

	}
	err = os.Mkdir(*folder_name, os.ModePerm)

	if err != nil {
		log.Fatal("Error with the given folder path")
	}
	var image_url string

	num_of_image_links := webhoundJson.MetaData.TotalItems

	for key := range num_of_image_links {
		image_url = webhoundJson.Data[key].Attributes.ImageUrl.Value
		image_id := webhoundJson.Data[key].Id
		new_image_name := *folder_name + "/" + image_id + ".jpg"

		image_file, err := os.Create(new_image_name)

		if err != nil {
			log.Fatal(err.Error())
		}

		defer image_file.Close()

		response, err := http.Get(image_url)

		_, err = io.Copy(image_file, response.Body)

		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Image download completed")

}
