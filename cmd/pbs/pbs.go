package main

import (
   "fmt"
   "github.com/89z/format/hls"
   "github.com/89z/mech/pbs"
   "net/http"
   "net/url"
   "os"
   "sort"
)

func doWidget(address string, bandwidth int, info bool) error {
   getter, err := pbs.NewWidgeter(address)
   if err != nil {
      return err
   }
   widget, err := getter.Widget()
   if err != nil {
      return err
   }
   addr := widget.HLS()
   fmt.Println("GET", addr)
   res, err := http.Get(addr)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   master, err := hls.NewMaster(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   fmt.Println(widget)
   sort.Sort(hls.Bandwidth{master, bandwidth})
   for _, stream := range master.Stream {
      if info {
         fmt.Println(stream)
      } else {
         audio := master.GetMedia(stream)
         if audio != nil {
            err := download(audio.URI, widget.Slug)
            if err != nil {
               return err
            }
         }
         return download(stream.URI, widget.Slug)
      }
   }
   return nil
}

func download(addr *url.URL, base string) error {
   fmt.Println("GET", addr)
   res, err := http.Get(addr.String())
   if err != nil {
      return err
   }
   seg, err := hls.NewSegment(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   if err := res.Body.Close(); err != nil {
      return err
   }
   file, err := os.Create(base + seg.Ext())
   if err != nil {
      return err
   }
   for i, info := range seg.Info {
      fmt.Print(seg.Progress(i))
      res, err := http.Get(info.URI.String())
      if err != nil {
         return err
      }
      if _, err := file.ReadFrom(res.Body); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return file.Close()
}
