package html

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ghetzel/friendscript/utils"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	defaults "github.com/mcuadros/go-defaults"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

type Document struct {
	document *goquery.Document
}

func (self *Document) String() string {
	if self.document == nil {
		return ``
	} else if out, err := self.document.Html(); err == nil {
		return out
	} else {
		return fmt.Sprintf("<!-- Error: %v -->", err)
	}
}

func (self *Document) Dump() (string, error) {
	if self.document == nil {
		return self.document.Html()
	} else {
		return ``, fmt.Errorf("No source document set")
	}
}

type Element struct {
	Namespace  string                 `json:"namespace,omitempty"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
	Text       string                 `json:"text"`
}

type Selection struct {
	Elements  []Element `json:"elements"`
	Length    int       `json:"length"`
	selection *goquery.Selection
}

type Commands struct {
	utils.Module
	scopeable utils.Scopeable
}

func New(scopeable utils.Scopeable) *Commands {
	cmd := &Commands{
		scopeable: scopeable,
	}

	cmd.Module = utils.NewDefaultExecutor(cmd)
	return cmd
}

type DumpArgs struct {
	// Whether to reformat the source being dumped.
	Format bool `json:"format" default:"true"`
}

// Dump the given source as a string.
func (self *Commands) Dump(source interface{}, args *DumpArgs) (string, error) {
	if args == nil {
		args = &DumpArgs{}
	}

	defaults.SetDefaults(args)

	if doc, err := self.parse(source); err == nil {
		if html, err := doc.Dump(); err == nil {
			var buf *bytes.Buffer

			if args.Format {
				buf = bytes.NewBuffer(nil)
				formatter := gohtml.NewWriter(buf)
				formatter.Write([]byte(html))
			} else {
				buf = bytes.NewBufferString(html)
			}

			return buf.String(), nil
		} else {
			return ``, err
		}
	} else {
		return ``, err
	}
}

type SelectArgs struct {
	// The source document to modify. Can be a URL, HTML string, byte array, open file handle, or Document object.
	Source interface{} `json:"source"`
}

// Retrieve HTML nodes matching the given CSS selector.
func (self *Commands) Find(selector string, args *SelectArgs) (*Selection, error) {
	if args == nil {
		return nil, fmt.Errorf("Source document must be provided.")
	}

	if sel, err := self.find(selector, args.Source); err == nil {
		return &Selection{
			Elements:  self.selectionToElements(sel),
			Length:    sel.Length(),
			selection: sel,
		}, nil
	} else {
		return nil, err
	}
}

// Remove the nodes matching the given Selection.
func (self *Commands) Remove(selector string, args *SelectArgs) error {
	if args == nil {
		return fmt.Errorf("Source document must be provided.")
	}

	if sel, err := self.find(selector, args.Source); err == nil {
		sel.Remove()
		return nil
	} else {
		return err
	}
}

type AttrArgs struct {
	// The source document to modify. Can be a URL, HTML string, byte array, open file handle, or Document object.
	Source interface{} `json:"source"`

	// The name of the attribute to set.
	Name string `json:"name"`

	// The value to set the attribute to.  If null, the attribute will be deleted.
	Value interface{} `json:"value"`
}

// Set the attribute to the given value for all elements matching the given
// selector.  If the value is null, the attribute will be removed.
func (self *Commands) Attr(selector string, args *AttrArgs) error {
	if args == nil {
		return fmt.Errorf("Source document must be provided.")
	} else if args.Name == `` {
		return fmt.Errorf("Must specify an attribute name to set.")
	}

	if sel, err := self.find(selector, args.Source); err == nil {
		if args.Value == nil {
			sel.RemoveAttr(args.Name)
		} else if value, err := stringutil.ToString(args.Value); err == nil {
			sel.SetAttr(args.Name, value)
		} else {
			return fmt.Errorf("value error: %v", err)
		}
	} else {
		return err
	}

	return nil
}

// Parse an HTML document and return an object that can be used with other
// functions in this module.
func (self *Commands) parse(source interface{}) (*Document, error) {
	doc := &Document{}
	var err error

	// attempt to parse the source document in some way
	if asDoc, ok := source.(*Document); ok {
		return asDoc, nil

	} else if asString, ok := source.(string); ok {
		if strings.HasPrefix(asString, `http`) {
			doc.document, err = goquery.NewDocument(asString)
		} else {
			doc.document, err = goquery.NewDocumentFromReader(
				strings.NewReader(asString),
			)
		}

	} else if asBytes, ok := source.([]byte); ok {
		doc.document, err = goquery.NewDocumentFromReader(
			bytes.NewBuffer(asBytes),
		)

	} else if asReader, ok := source.(io.Reader); ok {
		doc.document, err = goquery.NewDocumentFromReader(asReader)

	} else {
		return nil, fmt.Errorf("Source value could not be parsed.")

	}

	// ...then return it
	if err != nil {
		return nil, fmt.Errorf("error parsing source: %v", err)
	} else if doc.document == nil {
		return nil, fmt.Errorf("Source value could not be parsed.")
	} else {
		return doc, nil
	}
}

// Return the subselection for the given selector against the given source.
func (self *Commands) find(selector string, source interface{}) (*goquery.Selection, error) {
	if source == nil {
		return nil, fmt.Errorf("Source document must be provided.")
	}

	var gqSelection *goquery.Selection

	if fsSelection, ok := source.(*Selection); ok {
		gqSelection = fsSelection.selection
	} else if doc, err := self.parse(source); err == nil {
		gqSelection = doc.document.Selection
	} else {
		return nil, err
	}

	if gqSelection != nil {
		return gqSelection.Find(selector), nil
	} else {
		return nil, fmt.Errorf("Failed to locate source selection")
	}
}

// Take a selection and make it into Elements
func (self *Commands) selectionToElements(selection *goquery.Selection) []Element {
	elements := make([]Element, 0)

	if selection != nil {
		selection.Each(func(i int, s *goquery.Selection) {
			for _, node := range s.Nodes {
				switch node.Type {
				case html.ElementNode:
					element := Element{
						Namespace:  node.Namespace,
						Name:       node.Data,
						Attributes: make(map[string]interface{}),
						Text:       s.Text(),
					}

					for _, attr := range node.Attr {
						key := attr.Key

						if ns := attr.Namespace; ns != `` {
							key = ns + `:` + key
						}

						element.Attributes[key] = typeutil.Auto(attr.Val)
					}

					elements = append(elements, element)
				}
			}
		})
	}

	return elements
}
