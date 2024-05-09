package helper

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func fileNameToBytesArr(FileName string) []byte {
   file, err := os.Open(FileName)
   if err != nil {
	  fmt.Println(err)
	  return nil
   }
   defer file.Close()

   // Get the file size
   stat, err := file.Stat()
   if err != nil {
	  fmt.Println(err)
	  return nil
   }

   bs := make([]byte, stat.Size())
   _, err = bufio.NewReader(file).Read(bs)
   if err != nil && err != io.EOF {
	  fmt.Println(err)
	  return nil
   }
   return bs
}

/*
Given a byte array `bs`, decode it to a file with name `FileName`
*/
func bytesArrToFile(FileName string, bs []byte) {
   file, err := os.Create(FileName)
   if err != nil {
     fmt.Println(err)
     return
   }
   defer file.Close()


   // TODO: check if this is right
   _, err = file.Write(bs)
   if err != nil {
     fmt.Println(err)
     return
   }
}