package repositories

import (
	"fmt"
	"github.com/gjongenelen/redipo"
	"github.com/gjongenelen/redipo/cache"
	"github.com/google/uuid"
)

type Test struct {
	// Unique identifier for the test
	Id uuid.UUID `json:"id"`
	// Date on which test was taken
	Date int64 `json:"date"`
	// The object the test was performed on
	Object uuid.UUID `json:"object"`
	// Measurement results
	Values map[string]interface{} `json:"values"`
	// User notes
	Notes string `json:"notes"`
}

func NewTestsRepo(redis redipo.ManagerInterface) *TestsRepo {
	repo := redis.LoadDbRepo("tests")
	repo.SetFactory(func() interface{} {
		return &Test{}
	})
	repo.SetCaching(cache.NewMemoryCache())

	return &TestsRepo{
		repo,
	}
}

type TestsRepo struct {
	redipo.RepoInterface
}

func (t *TestsRepo) InitObjectIndex(object uuid.UUID) error {
	temp := uuid.New()

	err := t.AddToIndex(fmt.Sprintf("object_%s", object.String()), temp)
	if err != nil {
		return err
	}

	return t.RemoveFromIndex(fmt.Sprintf("object_%s", object.String()), temp)
}

func (t *TestsRepo) AddTestToObjectIndex(test Test) error {
	return t.AddToIndex(fmt.Sprintf("object_%s", test.Object.String()), test.Id)
}

func (t *TestsRepo) FindTestsByObjectId(object uuid.UUID) ([]*Test, error) {
	keys, err := t.GetIndex(fmt.Sprintf("object_%s", object.String()))
	if err != nil {
		return nil, err
	}

	tests := []*Test{}
	for _, id := range keys {
		test, err := t.Get(id)
		if err != nil {
			fmt.Printf("Test %s not found, originating from index %s", id, object)
			continue
		}
		if _, ok := test.(*Test); !ok {
			fmt.Printf("Test %s invalid, skipping", id)
			continue
		}
		tests = append(tests, test.(*Test))
	}

	return tests, nil
}
