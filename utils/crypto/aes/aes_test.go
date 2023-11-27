package aes

import (
	"log"
	"testing"

	
	"SuperpayProAdminWeb/utils/crypto/aes"
	
)

func TestBootConfig(t *testing.T) {
//	dec, _ := base64.Decode([]byte("pC/G7TjQmEqlhGI2FT0NNuqLC2L8HBmxRjLHfDLW2JCuJrIsvSQ26gn/T8m2yVhA"))
//	data, _ := Decrypt(dec, []byte("jfkgotmyvhspcandxlrwebquiz"))
 out,err:=aes.Base64Encrypt([]byte("9804659871*CUSTOMER"),[]byte(This is Very Nice Key))
	log.Println(string(out))
}
