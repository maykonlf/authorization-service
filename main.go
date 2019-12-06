package main

import (
	"fmt"
	"github.com/casbin/casbin"
)

func main() {
	e := casbin.NewEnforcer("rbac_model.conf", "roles.csv")

	e.AddPolicy("root", "*", "*", "*")

	e.AddPolicy("admin", "company1", "*", "*")
	e.AddPolicy("read_only", "company1", "*", "list")
	e.AddPolicy("read_only", "company1", "*", "get")

	//e.AddNamedPolicySafe("admin")

	//e.AddGroupingPolicy("maykon", "read_only", "company1")
	//e.AddGroupingPolicy("root", "root", "*")

	//fmt.Println(e.Enforce("maykon", "company1", "account", "list"))
	//fmt.Println(e.Enforce("maykon", "company1", "account", "delete"))
	//fmt.Println(e.Enforce("root", "company2", "account", "delete"))

	e.AddRoleForUserInDomain("root", "root", "*")

	e.SavePolicy()

	fmt.Println(e.Enforce("root", "company1", "account", "delete"))
	fmt.Println(e.Enforce("root", "company1", "account", "list"))
	fmt.Println(e.GetAllRoles())
}
