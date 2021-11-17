package book

import (
	"bookstop/graph/model"
	"strconv"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"google.golang.org/api/books/v1"
)

var converter = md.NewConverter("", true, nil)

func toMd(volume *books.Volume) *books.Volume {
	if volume.VolumeInfo.Description != "" {
		markdown, err := converter.ConvertString(volume.VolumeInfo.Description)
		if err == nil {
			volume.VolumeInfo.Description = markdown
		}
	}
	return volume
}

func ToGraph(volume *books.Volume) *model.Book {
	if volume == nil {
		return nil
	}

	if volume.VolumeInfo.PublishedDate == "" {
		// FIXME: temporary default to 0000
		volume.VolumeInfo.PublishedDate = "0000"
	}
	publishedYear, _ := strconv.Atoi(volume.VolumeInfo.PublishedDate[0:4])

	val := model.Book{
		ID:            volume.Id,
		Title:         volume.VolumeInfo.Title,
		Authors:       volume.VolumeInfo.Authors,
		Description:   volume.VolumeInfo.Description,
		PublishedYear: publishedYear,
	}

	if volume.VolumeInfo.Subtitle != "" {
		val.Subtitle = &volume.VolumeInfo.Subtitle
	}
	if volume.VolumeInfo.ImageLinks != nil && volume.VolumeInfo.ImageLinks.Thumbnail != "" {
		imageUrl := "https://books.google.com/books/publisher/content/images/frontcover/" + volume.Id + "?fife=w300"
		val.ImageURL = &imageUrl
	}
	if volume.VolumeInfo.IndustryIdentifiers != nil && len(volume.VolumeInfo.IndustryIdentifiers) > 0 {
		val.Isbn = &volume.VolumeInfo.IndustryIdentifiers[len(volume.VolumeInfo.IndustryIdentifiers)-1].Identifier
	}

	return &val
}
