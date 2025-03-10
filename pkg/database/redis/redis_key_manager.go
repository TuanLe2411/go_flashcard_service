package redis

func GetCategoriesKey(userId string) string {
	return "category:" + userId
}

func GetFlashcardsKey(userId string, categoryId string) string {
	return "flashcard:" + userId + ":" + categoryId
}
