package sync

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/api"
	"github.com/leonhfr/mochi/filesystem"
	"github.com/leonhfr/mochi/parser"
)

type requestKind int

const (
	createRequest requestKind = iota
	updateRequest
	archiveRequest
)

func (cr *CardResult) increment(kind requestKind) {
	switch kind {
	case createRequest:
		cr.Created++
	case updateRequest:
		cr.Updated++
	case archiveRequest:
		cr.Archived++
	}
}

type cardRequest struct {
	kind       requestKind
	id         string
	deckID     string
	content    string
	templateID string
	archived   bool
	fields     map[string]api.Field
	images     []syncImage
}

type syncImage struct {
	attachment api.Attachment
	path       string
	hash       string
}

func newCreateCardRequest(job *deckJob, card parser.Card, fs filesystem.Interface) (cardRequest, error) {
	content, fields := newCardContent(job, card)
	var images []syncImage
	for path, image := range card.Images {
		attachment, hash, err := newImageAttachment(path, image, fs)
		if err != nil {
			return cardRequest{}, err
		}
		if len(hash) > 0 {
			images = append(images, syncImage{
				attachment: attachment,
				path:       path,
				hash:       hash,
			})
		}
	}
	return cardRequest{
		kind:       createRequest,
		deckID:     job.id,
		content:    content,
		templateID: job.template.TemplateID,
		fields:     fields,
		images:     images,
	}, nil
}

func newUpdateCardRequest(job *deckJob, id string, card parser.Card, lock *Lock, fs filesystem.Interface) (cardRequest, error) {
	content, fields := newCardContent(job, card)
	var images []syncImage
	for path, image := range card.Images {
		attachment, hash, err := newImageAttachment(path, image, fs)
		if err != nil {
			return cardRequest{}, err
		}
		existingHash, ok := lock.getImageHash(id, path)
		if (!ok || existingHash != hash) && len(hash) > 0 {
			images = append(images, syncImage{
				attachment: attachment,
				path:       path,
				hash:       hash,
			})
		}
	}
	return cardRequest{
		kind:       updateRequest,
		id:         id,
		deckID:     job.id,
		content:    content,
		templateID: job.template.TemplateID,
		fields:     fields,
		images:     images,
	}, nil
}

func newArchiveCardRequest(id string) cardRequest {
	return cardRequest{
		id:       id,
		kind:     archiveRequest,
		archived: true,
	}
}

func processCardRequest(ctx context.Context, req cardRequest, lock *Lock, client Client) error {
	//nolint:prealloc
	var attachments []api.Attachment
	for _, image := range req.images {
		attachments = append(attachments, image.attachment)
	}

	switch req.kind {
	case createRequest:
		card, err := client.CreateCard(ctx, api.CreateCardRequest{
			DeckID:      req.deckID,
			Content:     req.content,
			TemplateID:  req.templateID,
			Fields:      req.fields,
			Attachments: attachments,
		})
		if err != nil {
			return err
		}
		setImages(card.ID, req.images, lock)
		return nil
	case updateRequest:
		card, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
			DeckID:      req.deckID,
			Content:     req.content,
			TemplateID:  req.templateID,
			Fields:      req.fields,
			Attachments: attachments,
		})
		if err != nil {
			return err
		}
		setImages(card.ID, req.images, lock)
		return nil
	case archiveRequest:
		_, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
			Archived: true,
		})
		return err
	default:
		return nil
	}
}

func newCardContent(job *deckJob, card parser.Card) (string, map[string]api.Field) {
	if !job.hasTemplate {
		return card.Content, nil
	}

	fields := make(map[string]api.Field)
	for id, field := range job.template.Fields {
		fields[id] = api.Field{
			ID:    id,
			Value: card.Fields[field],
		}
	}
	return "", fields
}

func setImages(id string, images []syncImage, lock *Lock) {
	for _, image := range images {
		lock.setImageHash(id, image.path, image.hash)
	}
}

func newImageAttachment(path string, image parser.Image, fs filesystem.Interface) (api.Attachment, string, error) {
	if !fs.FileExists(path) {
		return api.Attachment{}, "", nil
	}
	base64, hash, err := fs.Image(path)
	if err != nil {
		return api.Attachment{}, "", err
	}
	return api.Attachment{
		FileName:    fmt.Sprintf("%s.%s", image.FileName, image.Extension),
		ContentType: image.ContentType,
		Data:        string(base64),
	}, hash, nil
}
