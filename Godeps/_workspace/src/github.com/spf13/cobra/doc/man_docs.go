// Copyright 2015 Red Hat Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package doc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	mangen "github.com/cpuguy83/go-md2man/md2man"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gophergala2016/gophersiesta/Godeps/_workspace/src/github.com/spf13/pflag"
)

// GenManTree will generate a man page for this command and all decendants
// in the directory given. The header may be nil. This function may not work
// correctly if your command names have - in them. If you have `cmd` with two
// subcmds, `sub` and `sub-third`. And `sub` has a subcommand called `third`
// it is undefined which help output will be in the file `cmd-sub-third.1`.
func GenManTree(cmd *cobra.Command, header *GenManHeader, dir string) error {
	if header == nil {
		header = &GenManHeader{}
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsHelpCommand() {
			continue
		}
		if err := GenManTree(c, header, dir); err != nil {
			return err
		}
	}
	needToResetTitle := header.Title == ""

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".1"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := GenMan(cmd, header, f); err != nil {
		return err
	}

	if needToResetTitle {
		header.Title = ""
	}
	return nil
}

// GenManHeader is a lot like the .TH header at the start of man pages. These
// include the title, section, date, source, and manual. We will use the
// current time if Date if unset and will use "Auto generated by spf13/cobra"
// if the Source is unset.
type GenManHeader struct {
	Title   string
	Section string
	Date    *time.Time
	date    string
	Source  string
	Manual  string
}

// GenMan will generate a man page for the given command and write it to
// w. The header argument may be nil, however obviously w may not.
func GenMan(cmd *cobra.Command, header *GenManHeader, w io.Writer) error {
	if header == nil {
		header = &GenManHeader{}
	}
	b := genMan(cmd, header)
	final := mangen.Render(b)
	_, err := w.Write(final)
	return err
}

func fillHeader(header *GenManHeader, name string) {
	if header.Title == "" {
		header.Title = strings.ToUpper(strings.Replace(name, " ", "\\-", -1))
	}
	if header.Section == "" {
		header.Section = "1"
	}
	if header.Date == nil {
		now := time.Now()
		header.Date = &now
	}
	header.date = (*header.Date).Format("Jan 2006")
	if header.Source == "" {
		header.Source = "Auto generated by spf13/cobra"
	}
}

func manPreamble(out io.Writer, header *GenManHeader, name, short, long string) {
	dashName := strings.Replace(name, " ", "-", -1)
	fmt.Fprintf(out, `%% %s(%s)%s
%% %s
%% %s
# NAME
`, header.Title, header.Section, header.date, header.Source, header.Manual)
	fmt.Fprintf(out, "%s \\- %s\n\n", dashName, short)
	fmt.Fprintf(out, "# SYNOPSIS\n")
	fmt.Fprintf(out, "**%s** [OPTIONS]\n\n", name)
	fmt.Fprintf(out, "# DESCRIPTION\n")
	fmt.Fprintf(out, "%s\n\n", long)
}

func manPrintFlags(out io.Writer, flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if len(flag.Deprecated) > 0 || flag.Hidden {
			return
		}
		format := ""
		if len(flag.Shorthand) > 0 {
			format = "**-%s**, **--%s**"
		} else {
			format = "%s**--%s**"
		}
		if len(flag.NoOptDefVal) > 0 {
			format = format + "["
		}
		if flag.Value.Type() == "string" {
			// put quotes on the value
			format = format + "=%q"
		} else {
			format = format + "=%s"
		}
		if len(flag.NoOptDefVal) > 0 {
			format = format + "]"
		}
		format = format + "\n\t%s\n\n"
		fmt.Fprintf(out, format, flag.Shorthand, flag.Name, flag.DefValue, flag.Usage)
	})
}

func manPrintOptions(out io.Writer, command *cobra.Command) {
	flags := command.NonInheritedFlags()
	if flags.HasFlags() {
		fmt.Fprintf(out, "# OPTIONS\n")
		manPrintFlags(out, flags)
		fmt.Fprintf(out, "\n")
	}
	flags = command.InheritedFlags()
	if flags.HasFlags() {
		fmt.Fprintf(out, "# OPTIONS INHERITED FROM PARENT COMMANDS\n")
		manPrintFlags(out, flags)
		fmt.Fprintf(out, "\n")
	}
}

func genMan(cmd *cobra.Command, header *GenManHeader) []byte {
	// something like `rootcmd subcmd1 subcmd2`
	commandName := cmd.CommandPath()
	// something like `rootcmd-subcmd1-subcmd2`
	dashCommandName := strings.Replace(commandName, " ", "-", -1)

	fillHeader(header, commandName)

	buf := new(bytes.Buffer)

	short := cmd.Short
	long := cmd.Long
	if len(long) == 0 {
		long = short
	}

	manPreamble(buf, header, commandName, short, long)
	manPrintOptions(buf, cmd)
	if len(cmd.Example) > 0 {
		fmt.Fprintf(buf, "# EXAMPLE\n")
		fmt.Fprintf(buf, "```\n%s\n```\n", cmd.Example)
	}
	if hasSeeAlso(cmd) {
		fmt.Fprintf(buf, "# SEE ALSO\n")
		if cmd.HasParent() {
			parentPath := cmd.Parent().CommandPath()
			dashParentPath := strings.Replace(parentPath, " ", "-", -1)
			fmt.Fprintf(buf, "**%s(%s)**", dashParentPath, header.Section)
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}
		children := cmd.Commands()
		sort.Sort(byName(children))
		for i, c := range children {
			if !c.IsAvailableCommand() || c.IsHelpCommand() {
				continue
			}
			if cmd.HasParent() || i > 0 {
				fmt.Fprintf(buf, ", ")
			}
			fmt.Fprintf(buf, "**%s-%s(%s)**", dashCommandName, c.Name(), header.Section)
		}
		fmt.Fprintf(buf, "\n")
	}
	if !cmd.DisableAutoGenTag {
		fmt.Fprintf(buf, "# HISTORY\n%s Auto generated by spf13/cobra\n", header.Date.Format("2-Jan-2006"))
	}
	return buf.Bytes()
}
