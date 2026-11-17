package ldap

import "testing"

func TestDNIsActive(t *testing.T) {
	active := "CN=山田 太郎,OU=Proper,OU=Users,DC=example,DC=local"
	deleted := "CN=山田 太郎,OU=Delete Users,DC=example,DC=local"
	if !dnIsActive(active) {
		t.Fatal("active DN was treated as deleted")
	}
	if dnIsActive(deleted) {
		t.Fatal("Delete Users OU was treated as active")
	}
}
