package persist

import (
    "log"
    "os"
    "path"
    "io/ioutil"
    "regexp"
    "strings"
    
    "github.com/jmartin82/mmock/definition"
)

//FileBodyPersister persists body in file
type FileBodyPersister struct {
    PersistPath string
}

//Persist the body of the response to fiel if needed
func (fbp FileBodyPersister) Persist(per *definition.Persist, req *definition.Request, res *definition.Response) {
	if per.Name == ""{
        return
    }

    per.Name = fbp.replaceVars(req, per.Name)

    fileContent := []byte(res.Body)
    
    filePath := path.Join(fbp.PersistPath, per.Name);
    fileDir := path.Dir(filePath)

    err := os.MkdirAll(fileDir, 0644)
    if err != nil {
		panic(err)
	}

    err = ioutil.WriteFile(filePath, fileContent, 0644)

    if err != nil {
		panic(err)
	}
}

//NewFileBodyPersister creates a new FileBodyPersister
func NewFileBodyPersister(persistPath string) *FileBodyPersister {
	result := FileBodyPersister{PersistPath: persistPath}

    err := os.MkdirAll(result.PersistPath, 0644)
    if err != nil {
		panic(err)
	}

	return &result
}

func (fbp FileBodyPersister) replaceVars(req *definition.Request, input string) string {
	r := regexp.MustCompile(`\{\{\s*([^\}]+)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if tag == "request.body" {
			s = req.Body
			found = true
		} else if i := strings.Index(tag, "request.url."); i == 0 {
			s, found = req.GetURLPart(tag[12:], "Value")
		} else if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = req.GetQueryStringParam(tag[14:])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = req.GetCookieParam(tag[15:])
        }

		if !found {
			log.Printf("Defined tag {{%s}} not found\n", tag)
			return raw
		}
		return s
	})
}