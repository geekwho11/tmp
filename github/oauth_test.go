package github

import (
   "fmt"
   "testing"
   "time"
)

func TestOAuth(t *testing.T) {
   o, err := NewOAuth()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf(`1. Go to
%v

2. Enter this code
%v

3. Sign in to your GitHub Account
`, o.Verification_URI, o.User_Code)
   for range [9]struct{}{} {
      time.Sleep(9 * time.Second)
      x, err := o.Exchange()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", x)
      if x.Access_Token != "" {
         break
      }
   }
}