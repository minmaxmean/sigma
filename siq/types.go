package siq

// QuestionType represents well-known question types
const (
	QuestionTypeSimple  = "simple"
	QuestionTypeCat     = "cat"
	QuestionTypeAuction = "auction"
	QuestionTypeBagCat  = "bagCat"
	QuestionTypeSpider  = "spider"
	QuestionTypeSecret  = "secret"
	QuestionTypeNoRisk  = "noRisk"
	QuestionTypeSuper   = "super"
	QuestionTypeComplex = "complex"
	QuestionTypeMedia   = "media"
	QuestionTypeStake   = "stake"
	QuestionTypeFinal   = "final"
)

// ContentType represents content item types
const (
	ContentTypeText  = "text"
	ContentTypeImage = "image"
	ContentTypeAudio = "audio"
	ContentTypeVideo = "video"
	ContentTypeHtml  = "html"
)

// PlacementType represents content placement types
const (
	PlacementScreen     = "screen"
	PlacementReplic     = "replic"
	PlacementBackground = "background"
)

// ParamType represents parameter types
const (
	ParamTypeSimple    = "simple"
	ParamTypeContent   = "content"
	ParamTypeGroup     = "group"
	ParamTypeNumberSet = "numberSet"
)

// IsWellKnownType checks if a question type is well-known
func IsWellKnownType(questionType string) bool {
	wellKnownTypes := []string{
		QuestionTypeSimple,
		QuestionTypeCat,
		QuestionTypeAuction,
		QuestionTypeBagCat,
		QuestionTypeSpider,
		QuestionTypeSecret,
		QuestionTypeNoRisk,
		QuestionTypeSuper,
		QuestionTypeComplex,
		QuestionTypeMedia,
		QuestionTypeStake,
		QuestionTypeFinal,
	}

	for _, t := range wellKnownTypes {
		if t == questionType {
			return true
		}
	}
	return false
}
