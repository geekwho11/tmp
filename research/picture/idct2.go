package picture

func f2f(x float64) int32 {
   return int32(x * 4096 + 0.5)
}

func fsh(x int32) int32 {
   return x * 4096
}

func idct_1d(s0,s1,s2,s3,s4,s5,s6,s7,x0,x1,x2,x3 *int32) {
   p1, p4, p5 := new(int32), new(int32), new(int32)
   t0, t1, t2, t3 := new(int32), new(int32), new(int32), new(int32)
   p2 := s2
   p3 := s6
   *p1 = (*p2 + *p3) * f2f(0.5411961)
   *t2 = *p1 + *p3 * f2f(-1.847759065)
   *t3 = *p1 + *p2 * f2f( 0.765366865)
   p2 = s0
   p3 = s4
   *t0 = fsh(*p2 + *p3)
   *t1 = fsh(*p2 - *p3)
   *x0 = *t0 + *t3
   *x3 = *t0 - *t3
   *x1 = *t1 + *t2
   *x2 = *t1 - *t2
   t0 = s7
   t1 = s5
   t2 = s3
   t3 = s1
   *p3 = *t0 + *t2
   *p4 = *t1 + *t3
   *p1 = *t0 + *t3
   *p2 = *t1 + *t2
   *p5 = (*p3 + *p4) * f2f( 1.175875602)
   *t0 *= f2f(0.298631336)
   *t1 *= f2f(2.053119869)
   *t2 *= f2f(3.072711026)
   *t3 *= f2f(1.501321110)
   *p1 = *p5 + *p1 * f2f(-0.899976223)
   *p2 = *p5 + *p2 * f2f(-2.562915447)
   *p3 *= f2f(-1.961570560)
   *p4 *= f2f(-0.390180644)
   *t3 += *p1 + *p4
   *t2 += *p2 + *p3
   *t1 += *p2 + *p4
   *t0 += *p1 + *p3
}

func clamp(int64 x) int64 {
   if (x < 0) {
      return 0
   }
   if (x > 255) {
      return 255
   }
   return x
}

func idct_block(out []int64, stride int64, d [64]int64) {
   var v [64]int64
   for i := 0; i < 8; i++ {
      if d[8]==0 && d[16]==0 && d[24]==0 && d[32]==0 && d[40]==0 && d[48]==0 && d[56]==0 {
         v[0] = d[0] * 4
         v[8] = d[0] * 4
         v[16] = d[0] * 4
         v[24] = d[0] * 4
         v[32] = d[0] * 4
         v[40] = d[0] * 4
         v[48] = d[0] * 4
         v[56] = d[0] * 4
      } else {
         idct_1d(d[0],d[8],d[16],d[24],d[32],d[40],d[48],d[56])
         x0 += 512
         x1 += 512
         x2 += 512
         x3 += 512
         v[ 0] = (x0+t3) >> 10
         v[56] = (x0-t3) >> 10
         v[ 8] = (x1+t2) >> 10
         v[48] = (x1-t2) >> 10
         v[16] = (x2+t1) >> 10
         v[40] = (x2-t1) >> 10
         v[24] = (x3+t0) >> 10
         v[32] = (x3-t0) >> 10
      }
      d++
      v++
   }
   for i:=0, o:=out; i < 8; ++i,v+=8,o+=stride {
      idct_1d(v[0],v[1],v[2],v[3],v[4],v[5],v[6],v[7])
      x0 += 65536 + (128<<17)
      x1 += 65536 + (128<<17)
      x2 += 65536 + (128<<17)
      x3 += 65536 + (128<<17)
      o[0] = clamp((x0+t3) >> 17)
      o[7] = clamp((x0-t3) >> 17)
      o[1] = clamp((x1+t2) >> 17)
      o[6] = clamp((x1-t2) >> 17)
      o[2] = clamp((x2+t1) >> 17)
      o[5] = clamp((x2-t1) >> 17)
      o[3] = clamp((x3+t0) >> 17)
      o[4] = clamp((x3-t0) >> 17)
   }
}
