package cmd

import (
	"errors"
	"fmt"
	"io"
	"taothit/ggd/model"
)

func Generate(directiveSource *model.Directives, out *io.Writer, mode model.LogMode, srcPath string, rawSource string) error {
	if out == nil {
		return errors.New("datastructure output nil")
	}

	var sourceFile *model.SourceFile
	if srcPath != "" {
		sourceFile = directiveSource.CreateSourceFileFromPath(srcPath)
	} else if rawSource != "" {
		sourceFile = directiveSource.CreateSourceFileFromRawSource(rawSource)
	} else {
		return fmt.Errorf(`ggd: invalid directives (file="%s"; raw="%s"`, srcPath, rawSource)
	}

	ds := model.NewDatastructure(sourceFile, mode)
	if ds == nil {
		return fmt.Errorf("ggd: unknown datastructure for directives (%v)", directiveSource)
	}

	if err := ds.Print(out); err != nil {
		return fmt.Errorf("ggd: failed creating custom datastructure (%s): %v", ds, err)
	}
	return nil
}
