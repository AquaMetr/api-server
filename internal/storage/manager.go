package storage

import (
	"github.com/AquaMetr/api-server/internal/storage/repositories"
	"github.com/gjongenelen/redipo"
	"github.com/google/uuid"
)

type Storage struct {
	redipo.ManagerInterface

	objectsRepo *repositories.ObjectsRepo
	testsRepo   *repositories.TestsRepo
}

func NewStorage() *Storage {
	redis := redipo.New()

	return &Storage{
		ManagerInterface: redis,
		testsRepo:        repositories.NewTestsRepo(redis),
		objectsRepo:      repositories.NewObjectsRepo(redis),
	}
}

func (s *Storage) GetTestsRepo() *repositories.TestsRepo {
	return s.testsRepo
}

func (s *Storage) GetObjectsRepo() *repositories.ObjectsRepo {
	return s.objectsRepo
}

func (s *Storage) CreateObject(payload *repositories.Object) (*repositories.Object, error) {
	clone := *payload
	clone.Id = uuid.New()

	err := s.objectsRepo.Save(clone.Id, clone)
	if err != nil {
		return nil, err
	}

	err = s.testsRepo.InitObjectIndex(clone.Id)
	if err != nil {
		_ = s.objectsRepo.Delete(clone.Id)
		return nil, err
	}

	return &clone, nil
}

func (s *Storage) CreateTest(payload *repositories.Test) (*repositories.Test, error) {
	clone := *payload
	clone.Id = uuid.New()

	objectInterface, err := s.objectsRepo.Get(clone.Object)
	if err != nil {
		return nil, err
	}
	object := objectInterface.(*repositories.Object)

	err = s.testsRepo.Save(clone.Id, clone)
	if err != nil {
		return nil, err
	}

	err = s.testsRepo.AddTestToObjectIndex(clone)
	if err != nil {
		_ = s.testsRepo.Delete(clone.Id)
		return nil, err
	}

	if object.LatestTest < clone.Date {
		object.LatestTest = clone.Date

		_ = s.objectsRepo.Save(object.Id, object)
	}

	return &clone, nil
}
