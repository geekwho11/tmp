package main

import (
   "fmt"
   "github.com/89z/format/hls"
   "github.com/89z/mech/paramount"
   "net/http"
   "os"
   "sort"
)

func doManifest(guid string, bandwidth int, info bool) error {
   media, err := paramount.NewMedia(guid)
   if err != nil {
      return err
   }
   video, err := media.Video()
   if err != nil {
      return err
   }
   res, err := http.Get(video.Src)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   paramount.LogLevel.Dump(res.Request)
   master, err := hls.NewMaster(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   if info {
      fmt.Println(video.Title)
   }
   if bandwidth >= 1 {
      sort.Sort(hls.Bandwidth{master, bandwidth})
   }
   for _, stream := range master.Stream {
      if info {
         fmt.Println(stream)
      } else {
         return download(stream, video)
      }
   }
   return nil
}

func download(stream hls.Stream, video *paramount.Video) error {
   seg, err := newSegment(stream.URI.String())
   if err != nil {
      return err
   }
   res, err := http.Get(seg.Key.URI.String())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   paramount.LogLevel.Dump(res.Request)
   block, err := hls.NewCipher(res.Body)
   if err != nil {
      return err
   }
   file, err := os.Create(video.Base() + seg.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   for i, info := range seg.Info {
      fmt.Print(seg.Progress(i))
      res, err := http.Get(info.URI.String())
      if err != nil {
         return err
      }
      buf, err := block.Decrypt(info, res.Body)
      if err != nil {
         return err
      }
      if _, err := file.Write(buf); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return nil
}

func newSegment(addr string) (*hls.Segment, error) {
   res, err := http.Get(addr)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   paramount.LogLevel.Dump(res.Request)
   return hls.NewSegment(res.Request.URL, res.Body)
}
