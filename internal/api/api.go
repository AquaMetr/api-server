package api

import (
	"fmt"
	"github.com/AquaMetr/api-server/internal/storage"
	"github.com/AquaMetr/api-server/internal/storage/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func StartApi(storage *storage.Storage) error {
	api := echo.New()
	api.HideBanner = true
	api.HidePort = true

	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"healthy": true,
		})
	})

	v1Group := api.Group("/v1")

	v1Group.GET("/objects", func(c echo.Context) error {
		objects, err := storage.GetObjectsRepo().GetAllObjects()
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "could not fetch from storage",
			})
		}

		return c.JSON(http.StatusOK, objects)
	})

	v1Group.POST("/objects", func(c echo.Context) error {
		objectPayload := &repositories.Object{}
		err := c.Bind(objectPayload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		object, err := storage.CreateObject(objectPayload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "could not save to storage",
			})
		}

		return c.JSON(http.StatusOK, object)
	})

	v1Group.POST("/objects/:object/tests", func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("object"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "invalid object id",
			})
		}

		testPayload := &repositories.Test{}
		err = c.Bind(testPayload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		testPayload.Object = id

		test, err := storage.CreateTest(testPayload)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "could not save to storage",
			})
		}

		return c.JSON(http.StatusOK, test)
	})

	v1Group.GET("/tests/:object", func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("object"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "invalid object id",
			})
		}

		tests, err := storage.GetTestsRepo().FindTestsByObjectId(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "could not fetch from storage",
			})
		}

		return c.JSON(http.StatusOK, tests)
	})

	port := ":80"

	log.Infof("Starting api-server on %s", port)
	return api.Start(port)
}
