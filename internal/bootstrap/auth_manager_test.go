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
		t.Errorf("Error%s", err.Error())
	}
	auth := AuthManager{Tkn: &tkn, Cip: &ci, Kid: "presence",DurHours: 24}
	tk, err := auth.CreateToken(100, 120)
	if err != nil {
		t.Errorf("Error%s", err.Error())
	}
	session, err := auth.ValidateToken(tk)
	if err != nil {
		t.Errorf("Error%s", err.Error())
	}
	if session.SystemId != 100 {
		t.Errorf("Error SystemId%d", session.SystemId)
	}

	if session.Stub != 120 {
		t.Errorf("Error Stub%d", session.SystemId)
	}
}
