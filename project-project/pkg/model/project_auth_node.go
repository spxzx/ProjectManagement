package model

import "strings"

type ProjectAuthNode struct {
	Id   int64
	Auth int64
	Node string
}

type ProjectNodeAuthTree struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	Pnode    string
	Key      string
	Checked  bool
	Children []*ProjectNodeAuthTree
}

func ToAuthNodeTreeList(list []*ProjectNode, checkedList []string) []*ProjectNodeAuthTree {
	checkedMap := make(map[string]struct{})
	for _, v := range checkedList {
		checkedMap[v] = struct{}{}
	}
	var roots []*ProjectNodeAuthTree
	for _, v := range list {
		paths := strings.Split(v.Node, "/")
		if len(paths) == 1 {
			checked := false
			if _, ok := checkedMap[v.Node]; ok {
				checked = true
			}
			//根节点
			root := &ProjectNodeAuthTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeAuthTree{},
				Checked:  checked,
				Key:      v.Node,
			}
			roots = append(roots, root)
		}
	}
	for _, v := range roots {
		addAuthNodeChild(list, v, 2, checkedMap)
	}
	return roots
}

func addAuthNodeChild(list []*ProjectNode, root *ProjectNodeAuthTree, level int, checkedMap map[string]struct{}) {
	for _, v := range list {
		if strings.HasPrefix(v.Node, root.Node+"/") && len(strings.Split(v.Node, "/")) == level {
			//此根节点子节点
			checked := false
			if _, ok := checkedMap[v.Node]; ok {
				checked = true
			}

			child := &ProjectNodeAuthTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeAuthTree{},
				Checked:  checked,
				Key:      v.Node,
			}
			root.Children = append(root.Children, child)
		}
	}
	for _, v := range root.Children {
		addAuthNodeChild(list, v, level+1, checkedMap)
	}
}
