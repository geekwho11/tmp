package youtube

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

type Search struct {
   Contents struct {
      TwoColumnSearchResultsRenderer `json:"twoColumnSearchResultsRenderer"`
   }
}

func (r Request) NewSearch() (*Search, error) {
   buf := new(bytes.Buffer)
   json.NewEncoder(buf).Encode(r)
   req, err := http.NewRequest(
      "POST", "https://www.youtube.com/youtubei/v1/search", buf,
   )
   if err != nil {
      return nil, err
   }
   val := req.URL.Query()
   val.Set("key", "AIzaSyDCU8hByM-4DrUqRUYnGn-3llEO78bcxq8")
   req.URL.RawQuery = val.Encode()
   fmt.Println(invert, "POST", reset, req.URL)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("status %v", res.Status)
   }
   s := new(Search)
   json.NewDecoder(res.Body).Decode(s)
   return s, nil
}

func (s Search) VideoRenderers() []VideoRenderer {
   var vids []VideoRenderer
   for _, sect := range s.Contents.PrimaryContents.SectionListRenderer.Contents {
      for _, item := range sect.ItemSectionRenderer.Contents {
         vids = append(vids, item.VideoRenderer)
      }
   }
   return vids
}

type TwoColumnSearchResultsRenderer struct {
   PrimaryContents struct {
      SectionListRenderer struct {
         Contents []struct {
            ItemSectionRenderer struct {
               Contents []struct {
                  VideoRenderer `json:"videoRenderer"`
               }
            }
         }
      }
   }
}

type VideoRenderer struct {
   VideoID string
}
