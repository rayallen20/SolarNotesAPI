package biz

import (
	"cmp"
	"context"
	"errors"
	"slices"

	"SolarNotesAPI/model"
	"SolarNotesAPI/sys_error"
)

type Planet struct {
	ID    int64
	Name  string
	Title string
	Intro string
}

var planetDisplayOrder = []string{
	"Sun", "Mercury", "Venus", "Earth", "Mars",
	"Jupiter", "Saturn", "Uranus", "Neptune",
}

func planetRank(name string) int {
	if index := slices.Index(planetDisplayOrder, name); index >= 0 {
		return index
	}
	return len(planetDisplayOrder)
}

// List 将已创建的天体按前台展示顺序返回
func (Planet) List(ctx context.Context) ([]Planet, error) {
	if ctx == nil {
		return nil, &sys_error.InternalError{Cause: errors.New("biz.Planet.List 缺少请求上下文")}
	}

	var dao model.Planet
	rows, err := dao.FindAll(ctx)
	if err != nil {
		return nil, &sys_error.DatabaseError{Cause: err}
	}

	planets := make([]Planet, 0, len(rows))
	for _, row := range rows {
		planets = append(planets, Planet{
			ID:    row.ID,
			Name:  row.Name,
			Title: row.Title,
			Intro: row.Intro,
		})
	}

	// 未知名称的排序值相同。稳定排序会保留 DAO 提供的主键顺序。
	slices.SortStableFunc(planets, func(a, b Planet) int {
		return cmp.Compare(planetRank(a.Name), planetRank(b.Name))
	})
	return planets, nil
}
