package repositories

import (
	"github.com/gjongenelen/redipo"
	"github.com/google/uuid"
)

type ObjectType int64

const (
	Reserved ObjectType = iota
	FreshWaterAquarium
	SaltWaterAquarium
	Pond
	Pool
)

type Object struct {
	// Unique identifier for the object
	Id uuid.UUID `json:"id"`

	// User given name to the object
	Name string `json:"name"`

	// Type of the object
	Type ObjectType `json:"type"`

	// Date of latest performed test
	LatestTest int64 `json:"latest_test"`

	// Owner of the object
	Owner uuid.UUID `json:"owner"`
}

func NewObjectsRepo(redis redipo.ManagerInterface) *ObjectsRepo {
	repo := redis.LoadDbRepo("objects")
	repo.SetFactory(func() interface{} {
		return &Object{}
	})
	//repo.SetCaching(cache.NewMemoryCache())

	return &ObjectsRepo{
		repo,
	}
}

func (t *ObjectsRepo) GetAllObjects() ([]*Object, error) {
	results, err := t.GetAll()
	if err != nil {
		return nil, err
	}

	objects := []*Object{}
	for _, result := range results {
		objects = append(objects, result.(*Object))
	}

	return objects, nil
}

type ObjectsRepo struct {
	redipo.RepoInterface
}
