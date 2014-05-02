package gosplat

import (
	"strings"
	"text/template"
	"bytes"
	"encoding/json"
	 "io/ioutil"
	 "github.com/skratchdot/open-golang/open"
	 "fmt"
	 "os"
)


type Visualize struct {
	data [] interface{}
}

func NewVisualize() *Visualize {
	return &Visualize{make([]interface{}, 0)}
}

func (v *Visualize) Append(item interface{}) {
	v.data = append(v.data, item)
}

func (v *Visualize) plot(t string, options []map[string] interface{}) map[string] interface{} {
	if len(options) > 0 {
		return map[string]interface{}{"type": t, "data": v.data, "options": options[0]}
	}
	return map[string]interface{}{"type": t, "data": v.data, "options": make(map[string] interface{})}
}

func (v *Visualize) Linechart(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Linechart", options)
}

func (v *Visualize) Timeseries(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Timeseries", options)
}

func (v *Visualize) Table(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Table", options)
}
  
func (v *Visualize) Scatter(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Scatter", options)
}

func (v *Visualize) Candlestick(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Candlestick", options)
}

func (v *Visualize) Barchart(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Barchart", options)
}

func (v *Visualize) Columnchart(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Columnchart", options)
}

func (v *Visualize) Combochart(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Combochart", options)
}

func (v *Visualize) Treemap(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Treemap", options)
}

func (v *Visualize) Piechart(options ... map[string] interface{}) map[string] interface{}{
	return v.plot("Piechart", options)
}

type view struct {
	Title string
    Height string 
    Id string 
    Class string 
    Data string
    Options string 
	ColWidth string
}
  
func newView(title string, data map[string] interface{}, height string) (*view, error) {
	v := new(view)
	v.Title = title
	v.ColWidth = "grid_12"
	v.Height = height
	v.Id = strings.Replace(strings.ToLower(title), " ", "", -1)
	v.Class = data["type"].(string)
	d, err  := json.Marshal(data["data"].([]interface{}))
	if (err != nil) {
		return nil, err
	}
	v.Data = string(d)
	o, err := json.Marshal(data["options"].(map[string]interface{}))
	if (err != nil) {
		return nil, err
	}
	v.Options = string(o)
	return v, nil
}

type Row struct {
	Title string `json:"title"`
	Visualizations []*view	
}

func NewRow(title string) *Row {
	return &Row{title, make([]*view, 0)}
}
    
func (r *Row) Append(title string, viz map[string] interface{}, height ... string) error {
        if len(r.Visualizations) == 3 {
            return fmt.Errorf("Only up to 3 graphs per row for now")
		}
		if (len(height) > 0) {
			view, err := newView(title, viz, height[0]) 
			if (err != nil) {
				return err
			}
			r.Visualizations = append(r.Visualizations, view)
		} else {
			view, err := newView(title, viz, "375px")
			if (err != nil ) {
				return err
			}
			
			r.Visualizations = append(r.Visualizations, view)
		}
		for _, v := range r.Visualizations {
			v.ColWidth = r.cols()
		}
		return nil
}

func (r *Row) cols() (string) {
	switch len(r.Visualizations) {
	case 1:
		return "grid_12"
	case 2:
		return "grid_6"
	case 3:
		return "grid_4"
	default:
		return "grid_12"
	}
}     
			
type Frame struct {
	title string
	rows []*Row
}

func NewFrame(title string) *Frame{
	return &Frame{title, make([]*Row, 0)}
}

func (f *Frame) AppendRow(row *Row) {
	f.rows = append(f.rows, row)
}

func (f *Frame) Append(title string, v map[string]interface{}, height ... string) {
	row := NewRow("")
	if (len(height) > 0) {
		row.Append(title, v, height[0])
	} else {
		row.Append(title, v,"375px")
	}
	f.rows = append(f.rows, row)
}
	 
func (f *Frame) Html() (*bytes.Buffer, error) {
	t := template.New("go2splat")
	t, err := t.Parse(HTML)
	if (err != nil) {
		return nil, err
	}
	h := &html{f.title, f.rows}
	out := bytes.NewBuffer([]byte{})
	err = t.Execute(out, h)
	return out, err
}


//     def tmp_write_markup(self):
//         if not os.path.exists(os.path.join(os.getcwd(), "tmp")):
//             os.makedirs(os.path.join(os.getcwd(), "tmp"))
//         f = tempfile.NamedTemporaryFile(dir=os.path.join(os.getcwd(), "tmp"), suffix=".html", delete=None)
//         f.write(self.markup().__str__())
//         f.close()
//         return f

func (f *Frame) Preview() error {
	html, err := f.Html()
	
	if (err != nil) {
		return nil
	}
	
	p, err := ioutil.TempFile("", "go2splat.preview.")
	
	if err != nil {
		return err
	}
	_, err = p.Write(html.Bytes())
	if (err != nil) {
		return err
	}
	p.Close()
	name := fmt.Sprintf("%s.html", p.Name())
	os.Rename(p.Name(), name)
	open.Run(name)
	return nil
}

type html struct {
	Title string
	Rows []*Row
}
