package core

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func CollectFixtures(searchFixturesPaths []string) ([]Fixture, error) {
	fixturePaths, err := getFixturesPaths(searchFixturesPaths)
	if err != nil {
		return nil, err
	}

	fixtures := make([]Fixture, 0, len(fixturePaths))

	for _, path := range fixturePaths {
		fixture, err := loadFixture(path)
		if err != nil {
			return nil, err
		}
		fixtures = append(fixtures, fixture)
	}

	return fixtures, nil
}

// GetFixturesPaths - search in paths and returns yaml-files paths.
func getFixturesPaths(pathsToSearch []string) ([]string, error) {
	var result []string
	var err error

	for _, path := range pathsToSearch {
		if path == "" {
			continue
		}

		err = filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() && isYaml(entry.Name()) {
				result = append(result, path)
			}

			return nil
		})
	}

	return result, err
}

func load(reader io.Reader) (Fixture, error) {
	var eat RawFixture
	var fixture Fixture

	err := yaml.NewDecoder(reader).Decode(&eat)
	if err != nil {
		return fixture, err
	}

	err = move(&eat, &fixture)
	if err != nil {
		return fixture, err
	}

	return fixture, nil
}

// move - перекладывает прочитанные данные в DTO
// тут можно делать всю логику с перекидыванием данных.
func move(eat *RawFixture, fixture *Fixture) error {
	fixture.request.Method = eat.Request.Method
	fixture.request.Path = eat.Request.Path
	fixture.request.Query = eat.Request.Query
	fixture.request.Headers = eat.Request.Headers
	fixture.request.Body = eat.Request.Body

	fixture.response.Code = eat.Response.Code
	fixture.response.Headers = eat.Response.Headers
	fixture.response.Body = eat.Response.Body

	return nil
}

// loadFixture - парсит yaml-файл и возвращает фикстуру
// flpath - путь к конкретному yaml-файлу с фикстурой.
func loadFixture(flpath string) (Fixture, error) {
	file, err := os.Open(flpath)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("Close file %v error: %v", flpath, err)
		}
	}()
	if err != nil {
		return Fixture{}, err
	}

	result, err := load(file)
	if err != nil {
		return Fixture{}, err
	}

	return result, nil
}

// isYaml - смотрит на суффикс файла и возвращает true если это yaml-файл.
func isYaml(v string) bool {
	switch {
	case strings.HasSuffix(v, ".yaml"):
		return true
	case strings.HasSuffix(v, ".yml"):
		return true
	}

	return false
}
