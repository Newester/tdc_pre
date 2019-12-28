package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"runtime"
)

const file = "./test.gob"

type User struct {
	Name, Pass string
}

func main() {
	var datato = &User{"Donald", "DuckPass"}
	var datafrom = new(User)
	Getkeys()
	err := Save(file, datato)
	Check(err)
	err = Append(file)
	Check(err)
	RsaEncrypter(file)
	RsaDecrypter(file)
	err = Load(file, datafrom)
	Check(err)
	fmt.Println(datafrom)

}

func RsaEncrypter(path string) {
	fileInfo, _ := os.Stat(path)
	bufsize := fileInfo.Size()
	buf := make([]byte, bufsize)
	file, _ := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	file.Seek(0, os.SEEK_SET)
	file.Read(buf)
	file.Truncate(0)
	file.Close()
	ciphertext := RSA_encrypter("test_PublicKey.pem", buf)
	fmt.Println(ciphertext)
	file, _ = os.OpenFile(path, os.O_APPEND, os.ModeAppend)
	file.Write(ciphertext)
	file.Close()
}

func RsaDecrypter(path string) {
	fileInfo, _ := os.Stat(path)
	bufsize := fileInfo.Size()
	buf := make([]byte, bufsize)
	file, _ := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	file.Read(buf)
	result := RSA_decrypter("test_private.pem", buf)
	file.Truncate(0)
	file.Close()
	file, _ = os.OpenFile(path, os.O_APPEND, os.ModeAppend)
	file.Write(result)
	file.Close()
}

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

// Append UnusedBytes
func Append(path string) error {
	fileInfo, _ := os.Stat(path)
	used := fileInfo.Size()
	fmt.Printf("file size: %d\n", used)
	file, err := os.OpenFile(path, os.O_APPEND, os.ModeAppend)
	if err == nil {
		b := make([]byte, 128-used)

		_, err = file.Write(b)
		if err != nil {
			file.Close()
			log.Fatalf("write file error %v", err)
		}
	}
	file.Close()
	fileInfo, _ = os.Stat(path)
	fmt.Printf("file size: %d\n", fileInfo.Size())

	return err
}

// Check ...
func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}

func Getkeys() {
	//得到私钥
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	x509_Privatekey := x509.MarshalPKCS1PrivateKey(privateKey)
	//创建一个用来保存私钥的以.pem结尾的文件
	fp, _ := os.Create("test_private.pem")
	defer fp.Close()
	//将私钥字符串设置到pem格式块中
	pem_block := pem.Block{
		Type:  "test_privateKey",
		Bytes: x509_Privatekey,
	}
	//转码为pem并输出到文件中
	pem.Encode(fp, &pem_block)

	//处理公钥,公钥包含在私钥中
	publickKey := privateKey.PublicKey
	//接下来的处理方法同私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	x509_PublicKey, _ := x509.MarshalPKIXPublicKey(&publickKey)
	pem_PublickKey := pem.Block{
		Type:  "test_PublicKey",
		Bytes: x509_PublicKey,
	}
	file, _ := os.Create("test_PublicKey.pem")
	defer file.Close()
	//转码为pem并输出到文件中
	pem.Encode(file, &pem_PublickKey)

}

//使用公钥进行加密
func RSA_encrypter(path string, msg []byte) []byte {
	//首先从文件中提取公钥
	fp, _ := os.Open(path)
	defer fp.Close()
	//测量文件长度以便于保存
	fileinfo, _ := fp.Stat()
	buf := make([]byte, fileinfo.Size())
	fp.Read(buf)
	//下面的操作是与创建秘钥保存时相反的
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码,得到一个interface类型的pub
	pub, _ := x509.ParsePKIXPublicKey(block.Bytes)
	//加密操作,需要将接口类型的pub进行类型断言得到公钥类型
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), msg)
	fmt.Printf("len of msg: %d\n", len(msg))
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(msg)
	fmt.Println(cipherText)
	return cipherText
}

//使用私钥进行解密
func RSA_decrypter(path string, cipherText []byte) []byte {
	//同加密时，先将私钥从文件中取出，进行二次解码
	fp, _ := os.Open(path)
	defer fp.Close()
	fileinfo, _ := fp.Stat()
	buf := make([]byte, fileinfo.Size())
	fp.Read(buf)
	block, _ := pem.Decode(buf)
	PrivateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	//二次解码完毕，调用解密函数
	afterDecrypter, _ := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, cipherText)
	return afterDecrypter
}
