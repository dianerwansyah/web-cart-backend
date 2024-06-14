package helper

import (
	"net/http"
)

func GenericGetHandler(getDataFunc func() (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := getDataFunc()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, data)
	}
}

func ConvertToInterface[T any](getDataFunc func() ([]T, error)) func() (interface{}, error) {
	return func() (interface{}, error) {
		data, err := getDataFunc()
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}
