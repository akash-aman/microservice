// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package annotation

import (
	"entgo.io/contrib/entgql"
	"github.com/vektah/gqlparser/v2/ast"
)

func HasPermissions(permissions []string) entgql.Directive {
	children := make(ast.ChildValueList, 0, len(permissions))
	for _, p := range permissions {
		children = append(children, &ast.ChildValue{
			Value: &ast.Value{
				Raw:  p,
				Kind: ast.StringValue,
			},
		})
	}
	return entgql.NewDirective(
		"hasPermissions",
		&ast.Argument{
			Name: "permissions",
			Value: &ast.Value{
				Children: children,
				Kind:     ast.ListValue,
			},
		},
	)
}

func HasRole(permission string) entgql.Directive {
	return entgql.NewDirective(
		"hasRole",
		&ast.Argument{
			Name: "permission",
			Value: &ast.Value{
				Raw:  permission,
				Kind: ast.StringValue,
			},
		},
	)
}
