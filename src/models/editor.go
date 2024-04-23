package models

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type Editor struct {
	Content string
}

type EditorModel struct {
	Redis *redis.Client
}

func (em EditorModel) Get(editorId string) (Editor, error) {
	ctx := context.Background()

	val, err := em.Redis.Get(ctx, editorId).Result()

	if err != nil {
		return Editor{}, errors.New("editor page not found")
	}

	return Editor{Content: val}, nil
}

func (em EditorModel) Create(editorId string) {
	ctx := context.Background()

	err := em.Redis.Set(ctx, editorId, "", 0).Err()

	if err != nil {
		panic(err)
	}
}

func (em EditorModel) Save(editorId string, e Editor) {
	ctx := context.Background()

	err := em.Redis.Set(ctx, editorId, e.Content, 0).Err()

	if err != nil {
		panic(err)
	}
}

func (em EditorModel) All() []string {
	ctx := context.Background()
	var editorIds []string

	vals, _ := em.Redis.Do(ctx, "KEYS", "*").Result()

	for _, v := range vals.([]interface{}) {
		editorIds = append(editorIds, v.(string))
	}

	return editorIds
}
