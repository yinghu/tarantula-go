package bootstrap

import (
	"testing"

	"gameclustering.com/internal/util"
)

func TestAuthManager(t *testing.T) {
	tkn := util.JwtHMac{Alg: "SHS256"}
	tkn.HMac()

	ci := util.Cipher{Ksz: 32}
	err := ci.AesGcm()
	if err != nil {
		t.Errorf("Error%s",err.Error())	
	}
	auth := AuthManager{Tkn: &tkn, Cip: &ci, Kid: "presence"}
	tk,err := auth.CreateToken(100,120)
	if err != nil {
		t.Errorf("Error%s",err.Error())	
	}
	err = auth.ValidateToken(tk)
	if err != nil {
		t.Errorf("Error%s",err.Error())	
	}
}
