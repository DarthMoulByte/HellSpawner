package hsutil

import (
	"regexp"
	"strings"
)

type FileTreeNode struct {
	Name string
	Id   int
	Open bool

	IsFile   bool
	FullPath string

	Children []*FileTreeNode
}

// GetMpqPath build MpqPath from node
func (v *FileTreeNode) GetMpqPath() MpqPath {
	// Split the mpq filename from internal path
	// ex: d2data.mpq/data\global\items\flp2ax.dc6
	re := regexp.MustCompile(`[\\|/]`)
	pnames := re.Split(v.FullPath, 2)
	return MpqPath{pnames[0], pnames[1]}
}

func BuildTreeWalk(curnode *FileTreeNode, curpath []string, fullpath string, prevpaths string, id int) int {
	if len(curpath) == 0 {
		return id
	}

	// take the next bit off curpath
	var next string
	next, curpath = curpath[0], curpath[1:]
	prevpaths = prevpaths + "\\" + next

	// see if next already exists
	for _, node := range curnode.Children {
		if strings.ToLower(node.Name) == strings.ToLower(next) {
			return BuildTreeWalk(node, curpath, fullpath, prevpaths, id) // node already exists, keep walking
		}
	}

	// otherwise, add it
	isfile := len(curpath) == 0
	// find the index to add at
	// this logic ensures that dirs are on top of the list and files are on the bottom
	index := -1
	for i, node := range curnode.Children {
		if !isfile && node.IsFile || !isfile && node.Name > next {
			index = i
			break
		} else if isfile && node.IsFile && node.Name > next {
			index = i
			break
		}
	}

	newnode := &FileTreeNode{}
	if index == -1 {
		// if index is -1, it's a file or its a dir and we searched the whole list and found no files
		// so append it to the end
		curnode.Children = append(curnode.Children, newnode)
	} else {
		// insert the new node at a specific index
		curnode.Children = append(curnode.Children, nil)
		copy(curnode.Children[index+1:], curnode.Children[index:])
		curnode.Children[index] = newnode
	}
	newnode.Name = next
	newnode.IsFile = isfile
	newnode.Children = make([]*FileTreeNode, 0)
	id++
	newnode.Id = id
	if newnode.IsFile { // if it's a file, stop
		newnode.FullPath = fullpath
		return id
	} else { // otherwise, keep walking
		newnode.FullPath = prevpaths
		return BuildTreeWalk(newnode, curpath, fullpath, prevpaths, id)
	}
}

func BuildFileTreeFromFileList(paths []string) *FileTreeNode {
	root := &FileTreeNode{}
	root.Name = "root"
	root.Children = make([]*FileTreeNode, 0)

	id := 0
	for _, p := range paths {
		pnames := strings.Split(p, string("\\"))
		id = BuildTreeWalk(root, pnames, p, "", id)
	}

	return root
}
