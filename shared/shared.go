package shared

import (
	pb "github.com/sigmawq/grpc-service/grpc"
)

type DataEntry struct {
	Id         string `json:"_id"`
	Categories struct {
		Subcategory string `json:"subcategory"`
	} `json:"categories"`
	Title struct {
		Ro string `json:"ro"`
		Ru string `json:"ru"`
	} `json:"title"`
	Type   string  `json:"type"`
	Posted float64 `json:"posted"`
}

func (de *DataEntry) ToGrpcFormat() *pb.Data {
	return &pb.Data{
		Id:          de.Id,
		Subcategory: de.Categories.Subcategory,
		TitleRo:     de.Title.Ro,
		TitleRu:     de.Title.Ru,
		Type:        de.Type,
		Posted:      de.Posted,
	}
}

type DataEntryDatabase struct {
	Id          string  `json:"id"`
	Subcategory string  `json:"subcategory"`
	TitleRo     string  `json:"title_ro"`
	TitleRu     string  `json:"title_ru"`
	Type        string  `json:"type"`
	Posted      float64 `json:"posted"`
}

func NewDataEntryFromGrpcFormat(src *pb.Data) DataEntryDatabase {
	data := DataEntryDatabase{
		Id:          src.Id,
		Subcategory: src.Subcategory,
		TitleRo:     src.TitleRo,
		TitleRu:     src.TitleRu,
		Type:        src.Type,
		Posted:      src.Posted,
	}

	return data
}
