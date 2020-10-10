package templates

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"
)

/*
	only templates_cache or template_files are used depending if cache is on or not
*/
type Templates struct {
	cache           bool
	templates_dir   string
	func_map        template.FuncMap
	templates_cache map[string]*template.Template
	template_files  map[string][]string
}

func New(cache bool, templates_dir string, funcs template.FuncMap) (*Templates, error) {
	templates_cache := make(map[string]*template.Template)
	template_files := make(map[string][]string)
	templates := Templates{
		cache:           cache,
		templates_dir:   templates_dir,
		func_map:        funcs,
		templates_cache: templates_cache,
		template_files:  template_files,
	}

	return &templates, nil
}

func (self *Templates) AddTemplate(name string, files []string) error {
	/* files should be added with base last */
	if self.cache {
		templates, err := self.buildTemplate(files)
		if err != nil {
			return err
		}
		self.templates_cache[name] = templates
	} else {
		self.template_files[name] = files
	}

	return nil
}

func (self *Templates) ExecuteTemplate(name string, w io.Writer, data interface{}) error {
	if self.cache {
		template, ok := self.templates_cache[name]
		if !ok {
			return errors.New("template does not exist")
		}
		template.ExecuteTemplate(w, "__base", data)
	} else {
		files, ok := self.template_files[name]
		if !ok {
			return errors.New("template does not exist")
		}
		template, err := self.buildTemplate(files)
		if err != nil {
			return err
		}
		template.ExecuteTemplate(w, "__base", data)
	}

	return nil
}

func (self *Templates) buildTemplate(files []string) (*template.Template, error) {
	include_files_pattern := self.templates_dir + "/_[a-zA-Z0-9]*.html"
	tmpl := template.New("").Funcs(self.func_map)
	include_files, err := filepath.Glob(include_files_pattern)
	if len(include_files) > 0 && err == nil {
		tmpl, err = tmpl.ParseGlob(include_files_pattern)
		if err != nil {
			return nil, err
		}
	}
	new_files := make([]string, 0, len(files))
	for i := range files {
		new_files = append(new_files, self.templates_dir+"/"+files[i])
	}
	tmpl, err = tmpl.ParseFiles(new_files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
