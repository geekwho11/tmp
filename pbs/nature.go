package pbs

import (
   "github.com/89z/format/json"
   "net/http"
   "net/url"
)

type Nature map[string]struct {
   Video_Iframe string
}

func NewNature(addr string) (*Nature, error) {
   req, err := http.NewRequest("GET", addr, nil)
   if err != nil {
      return nil, err
   }
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var (
      nat = new(Nature)
      sep = []byte(`"full_length":`)
   )
   if err := json.Decode(res.Body, sep, nat); err != nil {
      return nil, err
   }
   return nat, nil
}

func (n Nature) Widget() (*Widget, error) {
   for _, val := range n {
      addr, err := url.Parse(val.Video_Iframe)
      if err != nil {
         return nil, err
      }
      addr.Scheme = "https"
      return NewWidget(addr)
   }
   return nil, notFound{"video_iframe"}
}
