package platform

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

func initTranslator() (*ut.UniversalTranslator, ut.Translator) {
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	return uni, trans
}
