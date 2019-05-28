package main

import (
	"os"
	"fmt"
	"github.com/valyala/fasttemplate"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func latestAnnotatedTag(repo *git.Repository) (*object.Tag, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var tag *object.Tag

	err = tags.ForEach(func (ref *plumbing.Reference) error {
		obj, err := repo.TagObject(ref.Hash())
		if err != nil {
			return nil
		}
		if tag == nil {
			tag = obj
			return nil
		}
		if obj.Tagger.When.After(tag.Tagger.When) {
			tag = obj
		}
		return nil
	})

	return tag, err
}

func latestTag(repo *git.Repository) (*object.Tag, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	tag, err := repo.TagObject(ref.Hash())
	if err == nil {
		return tag, nil
	}

	return latestAnnotatedTag(repo)
}

// Basic example of how to list tags.
func main() {
	path := os.Args[1]

	r, _ := git.PlainOpen(path)
	lt, _ := latestTag(r)

	template := "tag.name: {{tag.name}}\ntag.when: {{tag.when}}\nversion.major: {{version.major}}"
	t := fasttemplate.New(template, "{{", "}}")
	s := t.ExecuteString(map[string]interface{}{
		"tag.name" : lt.Name,
		"tag.when" : lt.Tagger.When.Format("2006-01-02T15:04:05-0700"),
	})
	fmt.Printf("%s", s)
}
