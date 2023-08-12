package model

import "strings"

type ProjectNode struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	CreateAt int64
}

type ProjectNodeTree struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	Pnode    string
	Children []*ProjectNodeTree
}

func ToNodeTreeList(list []*ProjectNode) []*ProjectNodeTree {
	var roots []*ProjectNodeTree
	for _, v := range list {
		paths := strings.Split(v.Node, "/")
		if len(paths) == 1 {
			//根节点
			root := &ProjectNodeTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeTree{},
			}
			roots = append(roots, root)
		}
	}
	for _, v := range roots {
		addChild(list, v, 2)
	}
	return roots
}

func addChild(list []*ProjectNode, root *ProjectNodeTree, level int) {
	for _, v := range list {
		if strings.HasPrefix(v.Node, root.Node+"/") && len(strings.Split(v.Node, "/")) == level {
			//此根节点子节点
			child := &ProjectNodeTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeTree{},
			}
			root.Children = append(root.Children, child)
		}
	}
	for _, v := range root.Children {
		addChild(list, v, level+1)
	}
}
