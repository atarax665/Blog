package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	enforcer, err := casbin.NewEnforcerSafe("./model.conf", "./policy.csv")
	if err != nil {
		panic(err)
	}

	// check permission
	// enforcer.Enforce("user", "resource", "permission") returns true if the request is allowed, else returns false
	res := enforcer.Enforce("abhinav", "data1", "read")
	fmt.Println(res) // true as abhinav is admin

	res = enforcer.Enforce("alice", "data2", "read")
	fmt.Println(res) // false as alice does not have read permission to data2

	// add permission
	// enforcer.AddPolicy("user", "resource", "permission")
	enforcer.AddPolicy("alice", "data2", "read")

	enforcer.SavePolicy()

	// check permission for alice again
	res = enforcer.Enforce("alice", "data2", "read")
	fmt.Println(res) // true as alice now has read permission to data2

	// get all permissions
	// enforcer.GetPermissionsForUser("user") returns all permissions for the user
	permissions := enforcer.GetPermissionsForUser("alice")
	fmt.Println(permissions) // [[alice data2 read]]

}
