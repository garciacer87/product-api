package api

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type productValidator struct {
	*validator.Validate
	t ut.Translator
}

func (v *productValidator) translate(err error) []string {
	var result []string
	fmtErrs := err.(validator.ValidationErrors).Translate(v.t)

	for _, e := range fmtErrs {
		result = append(result, e)
	}

	return result
}

func newValidator() *productValidator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	v := validator.New()
	en_translations.RegisterDefaultTranslations(v, trans)

	v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	v.RegisterTranslation("sku", trans, func(ut ut.Translator) error {
		return ut.Add("sku", "invalid SKU value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("sku", fe.Field())
		return t
	})

	v.RegisterTranslation("notblank", trans, func(ut ut.Translator) error {
		return ut.Add("notblank", "{0} is blank", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notblank", fe.Field())
		return t
	})

	v.RegisterTranslation("altimages", trans, func(ut ut.Translator) error {
		return ut.Add("altimages", "{0} has an invalid url value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("altimages", fe.Field())
		return t
	})

	v.RegisterTranslation("url", trans, func(ut ut.Translator) error {
		return ut.Add("url", "{0} is not a valid url value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url", fe.Field())
		return t
	})

	v.RegisterValidation("sku", func(fl validator.FieldLevel) bool {
		sku := fl.Field().String()
		return validateSKU(sku)
	})

	v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		v := fl.Field().String()
		return validateBlank(v)
	})

	v.RegisterValidation("altimages", func(fl validator.FieldLevel) bool {
		arr := fl.Field().Interface().([]string)
		return validateAltImages(arr)
	})

	return &productValidator{v, trans}
}

//Validates SKU values
func validateSKU(sku string) bool {
	//validates if sku starts with FAL-
	if !strings.HasPrefix(sku, "FAL-") {
		return false
	}

	id, err := strconv.Atoi(sku[4:])
	if err != nil {
		return false
	}

	//validates range of sku id
	if id < 1000000 || id > 99999999 {
		return false
	}

	return true
}

//Validates if the value is not blank. e.g.: "   " or ""
func validateBlank(v string) bool {
	return strings.TrimSpace(v) != ""
}

//Validates if the values of the string slice are valid URL
func validateAltImages(arr []string) bool {
	for _, imgURL := range arr {
		if _, err := url.ParseRequestURI(imgURL); err != nil {
			return false
		}
	}

	return true
}
