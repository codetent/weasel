package weasel

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateId() string {
	id, err := gonanoid.Generate("0123456789abcdef", 12)
	if err != nil {
		panic(err)
	}
	return id
}
