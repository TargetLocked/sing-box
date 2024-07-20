package build_shared

import (
	"github.com/sagernet/sing-box/common/badversion"
	"github.com/sagernet/sing/common/shell"
)

func ReadTag() (string, error) {
	currentTag, err := shell.Exec("make", "-s", "internaltag").ReadOutput()
	if err != nil {
		return currentTag, err
	}
	return currentTag, nil
}

func ReadTagVersion() (badversion.Version, error) {
	currentTag, err := ReadTag()
	if err != nil {
		return badversion.Version{}, err
	}
	version := badversion.Parse(currentTag)
	return version, nil
}
