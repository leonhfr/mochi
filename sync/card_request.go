package sync

import (
	"context"
	"fmt"
	"sort"

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
	filename   string
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

func newCreateCardRequest(job *deckJob, filename string, card parser.Card, fs filesystem.Interface) (cardRequest, error) {
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
		filename:   filename,
		deckID:     job.id,
		content:    content,
		templateID: job.template.TemplateID,
		fields:     fields,
		images:     images,
	}, nil
}

func newUpdateCardRequest(job *deckJob, id, filename string, card parser.Card, lock *Lock, fs filesystem.Interface) (cardRequest, error) {
	content, fields := newCardContent(job, card)
	paths := make([]string, 0, len(card.Images))

	var images []syncImage
	for path, image := range card.Images {
		paths = append(paths, path)
		attachment, hash, err := newImageAttachment(path, image, fs)
		if err != nil {
			return cardRequest{}, err
		}
		existingHash, ok := lock.getImageHash(job.id, id, path)
		if (!ok || existingHash != hash) && len(hash) > 0 {
			images = append(images, syncImage{
				attachment: attachment,
				path:       path,
				hash:       hash,
			})
		}
	}

	lock.cleanImages(job.id, id, paths)

	return cardRequest{
		kind:       updateRequest,
		filename:   filename,
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

func processCardRequest(ctx context.Context, req cardRequest, lock *Lock, client Client, logger Logger) error {
	//nolint:prealloc
	var attachments []api.Attachment
	for _, image := range req.images {
		attachments = append(attachments, image.attachment)
	}
	sort.Slice(attachments, func(i, j int) bool {
		return attachments[i].FileName < attachments[j].FileName
	})

	switch req.kind {
	case createRequest:
		if len(attachments) == 0 {
			_, err := createCard(ctx, req, lock, client, logger)
			return err
		}

		for index, attachment := range attachments {
			if index == 0 {
				card, err := createCardWithAttachment(ctx, req, attachment, lock, client, logger)
				if err != nil {
					return err
				}
				req.id = card.ID

				continue
			}

			if err := updateCardWithAttachment(ctx, req, attachment, lock, client, logger); err != nil {
				return err
			}
		}
		return nil
	case updateRequest:
		if len(attachments) == 0 {
			return updateCard(ctx, req, lock, client, logger)
		}

		for _, attachment := range attachments {
			if err := updateCardWithAttachment(ctx, req, attachment, lock, client, logger); err != nil {
				return err
			}
		}
		return nil
	case archiveRequest:
		return archiveCard(ctx, req, client, logger)
	default:
		return nil
	}
}

func createCard(ctx context.Context, req cardRequest, lock *Lock, client Client, logger Logger) (api.Card, error) {
	card, err := client.CreateCard(ctx, api.CreateCardRequest{
		DeckID:     req.deckID,
		Content:    req.content,
		TemplateID: req.templateID,
		Fields:     req.fields,
	})
	if err != nil {
		logger.Errorf("Card creation failed: %s...", substring(req.content, 100))
		return card, err
	}
	logger.Infof("Created card with id %s, deck id %s", card.ID, card.DeckID)

	if err := lock.setCard(req.deckID, card.ID, req.filename); err != nil {
		return card, err
	}

	if err := setImages(req.deckID, card.ID, req.images, lock); err != nil {
		return card, err
	}

	return card, nil
}

func createCardWithAttachment(ctx context.Context, req cardRequest, attachment api.Attachment, lock *Lock, client Client, logger Logger) (api.Card, error) {
	card, err := client.CreateCard(ctx, api.CreateCardRequest{
		DeckID:      req.deckID,
		Content:     req.content,
		TemplateID:  req.templateID,
		Fields:      req.fields,
		Attachments: []api.Attachment{attachment},
	})
	if err != nil {
		logger.Errorf("Card creation failed: %s...", substring(req.content, 100))
		return card, err
	}
	logger.Infof("Created card with id %s, deck id %s, attachment %s", card.ID, card.DeckID, attachment.FileName)

	if err := lock.setCard(req.deckID, card.ID, req.filename); err != nil {
		return card, err
	}

	if err := setImages(req.deckID, card.ID, req.images, lock); err != nil {
		return card, err
	}

	return card, nil
}

func updateCard(ctx context.Context, req cardRequest, lock *Lock, client Client, logger Logger) error {
	card, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
		DeckID:     req.deckID,
		Content:    req.content,
		TemplateID: req.templateID,
		Fields:     req.fields,
	})
	if err != nil {
		logger.Errorf("Card update failed (id: %s): %s...", req.id, substring(req.content, 100))
		return err
	}
	logger.Infof("Updated card with id %s, deck id %s", card.ID, card.DeckID)

	if err := lock.setCard(req.deckID, card.ID, req.filename); err != nil {
		return err
	}

	return setImages(req.deckID, card.ID, req.images, lock)
}

func updateCardWithAttachment(ctx context.Context, req cardRequest, attachment api.Attachment, lock *Lock, client Client, logger Logger) error {
	card, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
		DeckID:      req.deckID,
		Content:     req.content,
		TemplateID:  req.templateID,
		Fields:      req.fields,
		Attachments: []api.Attachment{attachment},
	})
	if err != nil {
		logger.Errorf("Card update with attachment failed (id: %s): %s...", req.id, substring(req.content, 100))
		return err
	}
	logger.Infof("Updated card with id %s, deck id %s, attachment %s", card.ID, card.DeckID, attachment.FileName)

	if err := lock.setCard(req.deckID, card.ID, req.filename); err != nil {
		return err
	}

	return setImages(req.deckID, card.ID, req.images, lock)
}

func archiveCard(ctx context.Context, req cardRequest, client Client, logger Logger) error {
	card, err := client.UpdateCard(ctx, req.id, api.UpdateCardRequest{
		Archived: true,
	})
	if err != nil {
		logger.Errorf("Card archive failed (id: %s)", req.id)
	}
	logger.Infof("Archived card with id %s, deck id %s", card.ID, card.DeckID)
	return nil
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

func setImages(deckID, cardID string, images []syncImage, lock *Lock) error {
	for _, image := range images {
		if err := lock.setImageHash(deckID, cardID, image.path, image.hash); err != nil {
			return err
		}
	}
	return nil
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

func substring(text string, length int) string {
	if len(text) < length {
		return text
	}
	return text[:length]
}
