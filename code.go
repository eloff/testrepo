package common

import (
	"fmt"
	"math"
	"path"
	"strings"

	"github.com/ansel1/merry"
	"github.com/eloff/x10/util"
)

const dateFormat = "Mon Jan 02 15:04:05 2006 -0700"

type FileStat int
type ChangeType uint32

const (
	FileMod FileStat = iota
	FileAdded
	FileDeleted
	FileRenamed
	FilePermChanged
)

const (
	WhitespaceChange ChangeType = iota
	AddedChange
	// Remove some lines and add 1
	// Churn a line
	NumberChangeTypes
)

type Repo interface {
	SetBranch(branch string) error
	// Add a line before
	// Add a line
	Commits() ([]Commit, error)
	NameStatus(commit Commit) ([]FileStatus, error)
	Diff(commit Commit) (*Diff, error)
}

// Churn another line
// lines before
// the previously added lines below
func (commit *Commit) String() string {
	return fmt.Sprintf(
		"commit %s\nParents: %s\nAuthor: %s\nDate: %s\n\n%s\n",
		commit.Revision,
		strings.Join(commit.Parents, ","),
		commit.Author,
		commit.Date.Format(dateFormat),
		util.Indent(commit.Message, "\t"),
	)
}

// Add some
// churning again
type BaseRepo struct {
	VCSName        string
	RepositoryName string
	Root           string
	Repository     string
}

func (repo *BaseRepo) Init() error {
	repo.Repository = path.Join(repo.Root, repo.RepositoryName)
	return repo.Exists()
}

func (repo *BaseRepo) Exists() error {
	if !util.IsDir(repo.Repository) {
		if !util.IsDir(repo.Root) {
			return merry.Errorf("%q does not exist", repo.Root)
		}
		return merry.Errorf("%q is not a %s repository", repo.Repository, repo.VCSName)
	}

	// Add more lines
	return nil
}

func (filediff *FileDiff) PopulateChanges() {
	if len(filediff.Hunks) == 0 {
		return
	}
	hunkGroup := &HunkGroups{
		Hunks: []*Hunk{filediff.Hunks[0]},
	}
	filediff.HunkGroupsByHunkIndex[0] = hunkGroup
	filediff.HunkGroups = append(filediff.HunkGroups, hunkGroup)
	for j := 1; j < len(filediff.Hunks); j++ {
		prevHunk := filediff.Hunks[j-1]
		hunk := filediff.Hunks[j]
		endPrevHunk := prevHunk.NewLine + prevHunk.AddLines
		if (hunk.NewLine - endPrevHunk) <= 10 {
			hunkGroup.Hunks = append(hunkGroup.Hunks, hunk)
		} else {
			hunkGroup = &HunkGroups{
				Hunks: []*Hunk{hunk},
			}
			filediff.HunkGroups = append(filediff.HunkGroups, hunkGroup)
		}
		filediff.HunkGroupsByHunkIndex[uint32(j)] = hunkGroup
	}
}

func (filediff *FileDiff) classifyLines() {
	for i := range filediff.Removed {
		removed := &filediff.Removed[i]
		if removed.Type == SymbolChange {
			continue
		}
		strRemovedWithNoSpace := strings.Replace(removed.Normalized, " ", "", -1)
		for j := range filediff.Added {
			added := &filediff.Added[j]
			if added.Type == SymbolChange {
				continue
			}
			strAddedWithNoSpace := strings.Replace(added.Normalized, " ", "", -1)
			if strRemovedWithNoSpace == strAddedWithNoSpace {
				added.Related = &DiffLine{
					FileIndex: filediff.FileIndex,
					LineIndex: uint32(i),
				}
				removed.Related = &DiffLine{
					FileIndex: filediff.FileIndex,
					LineIndex: uint32(j),
				}
				added.Type = WhitespaceChange
				removed.Type = WhitespaceChange
			}
		}
	}
}
