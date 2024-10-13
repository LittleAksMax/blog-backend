package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func validateRequest(ctx *gin.Context, reqPtr interface{}) (map[string]string, bool) {
	newCtx := ctx.Copy()
	if err := newCtx.ShouldBindJSON(reqPtr); err != nil {
		// If validation fails, respond with an exrror
		if errs, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldError := range errs {
				fieldName := fieldError.Field()
				tag := fieldError.Tag()

				var key string
				// if we failed on the required tag
				if tag == "required" {
					key = fieldName
				} else {
					field, _ := reflect.TypeOf(reqPtr).Elem().FieldByName(fieldName)

					if fieldJSONName, ok := field.Tag.Lookup("json"); ok {
						key = fieldJSONName
					} else {
						key = fieldName
					}
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
