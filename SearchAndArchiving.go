package main
import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "crypto/sha256"
    "io"
    "log"
    "encoding/hex"
    "archive/zip"
    "errors"
    "bytes"
)

func main() {
  var Name string
  fmt.Println("Название файла для архивирования")
  fmt.Scan(&Name)

  FilePath, err := SearchFile(ConcatStrings("/", Name))
    if err != nil{
      log.Fatal(err)
    }

  HashName := HashFile(FilePath)
  HashName = ConcatStrings(HashName, ".zip")
  files := []string{FilePath}
    if err := ZipFiles(HashName+"zip", files); err != nil {
        panic(err)
    }
    fmt.Println("Zipped File:", HashName)

}

func ConcatStrings(LeftPart string, RightPart string) string {
  strings := []string{LeftPart, RightPart}
  buffer := bytes.Buffer{}
  for _, val := range strings {
    buffer.WriteString(val)
  }
  return buffer.String()
}

func SearchFile(FileName string) (string, error){
  var flag bool = false
  HomeDir := os.UserHomeDir()
  filepath.Walk(HomeDir, func(path string, info os.FileInfo, err error) error {
      if strings.Contains(path, FileName){
        FileName = path
        flag = true
      }
      return nil
  })
  if flag == false{
    return " ", errors.New("File not found")
  }
  return FileName, nil
}

func HashFile(FilePath string ) string{
  f, err := os.Open(FilePath)
 if err != nil {
   log.Fatal(err)
 }
 defer f.Close()

 h := sha256.New()
 if _, err := io.Copy(h, f); err != nil {
   log.Fatal(err)
 }

 hx := hex.EncodeToString(h.Sum(nil))
 return hx
}

func ZipFiles(filename string, files []string) error {

    newZipFile, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer newZipFile.Close()

    zipWriter := zip.NewWriter(newZipFile)
    defer zipWriter.Close()

    for _, file := range files {
        if err = AddFileToZip(zipWriter, file); err != nil {
            return err
        }
    }
    return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

    fileToZip, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer fileToZip.Close()

    info, err := fileToZip.Stat()
    if err != nil {
        return err
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
        return err
    }

    header.Name = filename

    header.Method = zip.Deflate

    writer, err := zipWriter.CreateHeader(header)
    if err != nil {
        return err
    }
    _, err = io.Copy(writer, fileToZip)
    return err
}
