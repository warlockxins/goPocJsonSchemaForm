package handler

import (
	"fmt"
	"os"
)

func GetScreens() ([]string, error) {
	files, err := os.ReadDir("./screens")
	if err != nil {
		fmt.Println("Error getting screen list", err)
		return nil, err
	}

	var directories []string = []string{}
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}

func ScreensToHandlerScreeenList(directories *[]string, current *string) []Screen {
	screens := []Screen{}

	for _, dir := range *directories {
		screens = append(screens, Screen{
			Name:    dir,
			Current: current != nil && dir == *current,
		})

	}

	return screens
}
