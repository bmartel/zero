# Zero

A light convenience wrapper on go-playground/validator to allow readable error messages to be generated.

A simple example of using zero validation
```go
type Post struct {
	zero.Validation
	Title string `valid:"min=3,max=64"`
	Body  string `valid:"min=10,max=512"`
}

func main() {
  // Assume this post was constructed from values of a web request
  post := Post{
    Title: "test",
    Body:  "this is the body content",
  }

  // Create a new validator instance with the target struct tag
  v := zero.New("valid")

  // Run the validator on the post struct
  msgs, valid := v.Validate(post)
  if !valid {
    // Any associated error messages will be populated
    log.Printf("%+v", msgs)
  }
}
```

Alternate error messages to be returned on validation failure
```go
func main() {

  customMessages := map[string]string{
    "lte": "%s cannot be greater than %s",
    "min": "%s cannot have a size less than %s",
  }

  v := zero.New("valid")

  // Will override all defaults, if some values are missing no return values will be made for those not found
  v.SetMessages(customMessages)

  // Individual messages can be set or overridden at a time
  v.SetMessage("hexcolor", "%s must represent a hexadecimal color")
}
```

Custom error messages per struct type under validation
```go
type Comment struct {
  zero.Validation
  Body string `valid:"min=10,max=512"`
}

type Post struct {
	Title string `valid:"min=3,max=64"`
	Body  string `valid:"min=10,max=512"`
}

// If a struct does not compose a zero.Validation struct, the Validates() func must be implemented
func (Post) Validates() map[string]string {
  return map[string]string{
    "title.min": "%s cannot have a length less than %s",
    "title.max": "title cannot have a length greater than 512",
    "body.min": "body cannot have less than 3 characters",
    "body.max": "%s cannot have more than %s characters",
  }
}

func main() {
  post := Post{
    Title: "t",
    Body:  "b",
  }

  comment := Comment{
    Body:  "b",
  }

  v := zero.New("valid")

  msgs, valid := v.Validate(post)
  if !valid {
    // Any error messages will first look to use the values returned from the Validates() func specified on the Post struct
    log.Printf("%+v", msgs)
  }

  msgs, valid = v.Validate(comment)
  if !valid {
    // Since the Comment struct uses zero.Validation it will use the default messages for its response
    log.Printf("%+v", msgs)
  }
}
```

Adding custom validators
```go
var isObjectIDHex = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
var containsHTMLContent = regexp.MustCompile(`</?\w+((\s+\w+(\s*=\s*(?:".*?"|'.*?'|[\^'">\s]+))?)+\s*|\s*)/?>`)

func isObjectId(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
	return isObjectIDHex.MatchString(field.String())
}

func containsHTML(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
	return containsHTMLContent.MatchString(field.String())
}

type Comment struct {
  zero.Validation
  ID   string `valid:"objectid"`
  Body string `valid:"html"`
}

func main() {

  comment := Comment{
    Body:  "b",
  }

  v := zero.New("valid")

  // Can add custom validators one at a time ...
  v.AddValidator("objectid", isObjectId, "%s must be a valid objectid")
  v.AddValidator("html", containsHTML, "%s must contain valid html")

  // ... or all at once!
  // v.AddValidators(map[string]zero.ValidatorFunc{
  //   "objectid": zero.ValidatorFunc{Message: "%s must be a valid objectid", Func: isObjectId},
  //   "html":     zero.ValidatorFunc{Message: "%s must contain valid html", Func: containsHTML},
  // })

  msgs, valid := v.Validate(comment)
  if !valid {
    // Will show the invalid messages registered for the html, and objectid validators
    log.Printf("%+v", msgs)
  }

}
```
