package youtube

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "mime"
   "net/http"
   "strings"
   "time"
)

func (f Format) Format(s fmt.State, verb rune) {
   if f.QualityLabel != "" {
      fmt.Fprint(s, "Quality:", f.QualityLabel)
   } else {
      fmt.Fprint(s, "Quality:", f.AudioQuality)
   }
   fmt.Fprint(s, " Bitrate:", f.Bitrate)
   if f.ContentLength >= 1 { // Tq92D6wQ1mg
      fmt.Fprint(s, " Size:", f.ContentLength)
   }
   fmt.Fprint(s, " Type:", f.MimeType)
   if verb == 'a' {
      fmt.Fprint(s, " URL:", f.URL)
   }
}

func (f Format) Write(dst io.Writer) error {
   req, err := http.NewRequest("GET", f.URL, nil)
   if err != nil {
      return err
   }
   LogLevel.Dump(req)
   var (
      begin = time.Now()
      content int64
   )
   for content < f.ContentLength {
      req.Header.Set(
         "Range", fmt.Sprintf("bytes=%v-%v", content, content+partLength-1),
      )
      end := time.Since(begin).Seconds()
      if end >= 1 {
         fmt.Print(format.Percent(content, f.ContentLength), "\t")
         fmt.Print(format.LabelSize(content), "\t")
         fmt.Println(format.LabelRate(float64(content)/end))
      }
      // this sometimes redirects, so cannot use http.Transport
      res, err := new(http.Client).Do(req)
      if err != nil {
         return err
      }
      if _, err := io.Copy(dst, res.Body); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
      content += partLength
   }
   return nil
}

const partLength = 10_000_000

// averageBitrate is not always available:
// Tq92D6wQ1mg
type Format struct {
   AudioQuality string
   Bitrate int
   ContentLength int64 `json:"contentLength,string"`
   Height int
   MimeType string
   QualityLabel string
   URL string
   Width int
}

func (f Format) Ext() (string, error) {
   exts, err := mime.ExtensionsByType(f.MimeType)
   if err != nil {
      return "", err
   }
   for _, ext := range exts {
      return ext, nil
   }
   return "", notFound{f.MimeType}
}

type Formats []Format

func (f Formats) Len() int {
   return len(f)
}

func (f Formats) MediaType() error {
   for i, form := range f {
      typ, param, err := mime.ParseMediaType(form.MimeType)
      if err != nil {
         return err
      }
      param["codecs"], _, _ = strings.Cut(param["codecs"], ".")
      f[i].MimeType = mime.FormatMediaType(typ, param)
   }
   return nil
}

func (f Formats) Swap(i, j int) {
   f[i], f[j] = f[j], f[i]
}

// We cannot use bitrate to sort, as you end up with different heights:
//
// ID          | 480    | 720 low | 720 high | 1080
// ------------|--------|---------|----------|-----
// 7WTEB7Qbt4U | 285106 | 286687  | 513601   | 513675
// RPjE9riEhtA | 584072 | 1169166 | 1693812  | 2151670
type Height struct {
   Formats
   Target int
}

func (h Height) Less(i, j int) bool {
   return h.distance(i) < h.distance(j)
}

func (h Height) distance(i int) int {
   diff := h.Formats[i].Height - h.Target
   if diff >= 0 {
      return diff
   }
   return -diff
}

type notFound struct {
   value string
}

func (n notFound) Error() string {
   return fmt.Sprintf("%q is not found", n.value)
}
