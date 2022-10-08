package nntp

import (
	"encoding/json"
	"fmt"
	"newsmere/internal/storage"
	"newsmere/internal/types"

	"gorm.io/gorm/clause"
)

func New(config json.RawMessage) (*Backend, error) {
	backend := new(Backend)
	err := json.Unmarshal(config, &backend)
	return backend, err
}

func (Backend) Type() string {
	return Type
}

func (b *Backend) Start() error {
	fmt.Printf("[Backend] %s-%s starting\n", b.Type(), b.Name)
	err := b.syncSubs()
	fmt.Printf("[Backend] %s-%s sync subscriptions finished\n",
		b.Type(), b.Name)
	return err
}

func (b *Backend) Stop() error {
	if err := b.client.Close(); err != nil {
		return err
	}
	b.client = nil

	return nil
}

func (b *Backend) Restart() error {
	if err := b.Stop(); err != nil {
		return err
	}

	if err := b.Start(); err != nil {
		return err
	}

	return nil
}

func (b *Backend) Status() string {
	if b.client == nil {
		return types.StatusDown
	}

	return types.StatusUp
}

func (b *Backend) ensureClient() error {
	if b.client != nil {
		return nil
	}

	addr := fmt.Sprintf("%s:%d", b.Server, b.Port)
	if b.Port == 0 {
		addr = b.Server + ":119"
	}

	var err error
	b.client, err = NewClient("tcp", addr)
	if err != nil {
		return err
	}

	if b.User != "" {
		_, err := b.client.Authenticate(b.User, b.Pass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Backend) syncSubs() error {
	if err := b.ensureClient(); err != nil {
		return err
	}

	db := storage.GetDb()

	// check subscriptions exist
	var count int64
	db.Model(&storage.Subscription{}).Where("source = ?", b.Name).Count(&count)

	var subs []storage.Subscription
	if count == 0 {
		groups, err := b.client.List("")
		if err != nil {
			return err
		}

		for _, g := range groups {
			subs = append(subs, storage.Subscription{
				Name:        g.Name,
				Description: g.Name,
				High:        int(g.High),
				Low:         int(g.Low),
				Source:      b.Name,
			})
		}

		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "source"}},
			DoUpdates: clause.AssignmentColumns([]string{"high", "low"}),
		}).Create(&subs)

		if result.Error != nil {
			return result.Error
		}

		subs = subs[:10]
	} else {
		result := db.Where("source = ?", b.Name).Limit(10).Find(&subs)
		if result.Error != nil {
			return result.Error
		}
	}

	// create groups
	var groups []storage.Group
	for _, s := range subs {
		groups = append(groups, storage.Group{
			Name:        s.Name,
			Description: s.Description,
			Source:      s.Source,
			Low:         s.Low,
			High:        s.High,
		})
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "Source"}},
		DoUpdates: clause.AssignmentColumns([]string{"high", "low"}),
	}).Create(&groups)

	return result.Error
}
