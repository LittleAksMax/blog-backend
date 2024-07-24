package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func validateRequest(ctx *gin.Context, reqPtr interface{}) (map[string]string, bool) {
	newCtx := ctx.Copy()
	if err := newCtx.ShouldBindJSON(reqPtr); err != nil {
		// If validation fails, respond with an error
		if errs, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldError := range errs {
				fieldName := fieldError.Field()
				field, _ := reflect.TypeOf(reqPtr).Elem().FieldByName(fieldName)

				fieldJSONName, ok := field.Tag.Lookup("json")

				var key string
				if ok {
					key = fieldJSONName
				} else {
					key = fieldName
				}
				errorMessages[key] = fieldError.Tag()
			}
			return errorMessages, false
		}

		// existence of "_" signifies a different error type
		return map[string]string{"_": err.Error()}, false
	}

	return nil, true
}
