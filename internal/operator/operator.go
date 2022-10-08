package operator

import (
	nntp_sv "newsmere/internal/service/nntp"
	"newsmere/internal/storage"
	"strings"
)

func New() *Operator {
	return &Operator{
		authorized: false,
	}
}

func (o *Operator) ListGroups(max int) ([]*storage.Group, error) {
	db := storage.GetDb()

	var groups []*storage.Group
	if max < 0 {
		result := db.Find(&groups)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := db.Limit(max).Find(&groups)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return groups, nil
}

func (o *Operator) GetGroup(name string) (*storage.Group, error) {
	if name == "" {
		return nil, nntp_sv.ErrNoSuchGroup
	}

	parts := strings.SplitN(name, ".", 2)

	if len(parts) < 2 {
		return nil, nntp_sv.ErrNoSuchGroup
	}

	db := storage.GetDb()

	var group *storage.Group
	result := db.Where("name = ? AND source = ?", parts[1], parts[0]).First(&group)
	if result.Error != nil {
		return nil, result.Error
	}

	return group, nil
}

func (o *Operator) GetArticle(group *nntp_sv.Group, id string) (
	*nntp_sv.Article, error) {
	return nil, nil
}

func (o *Operator) GetArticles(group *nntp_sv.Group, from, to int64) (
	[]nntp_sv.NumberedArticle, error) {
	return nil, nil
}

func (o *Operator) Authorized() bool {
	return o.authorized
}

func (o *Operator) Authenticate(user, pass string) (nntp_sv.Operator, error) {
	o.authorized = true
	return o, nil
}
