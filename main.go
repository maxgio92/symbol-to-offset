package main

import (
	"debug/elf"
	"errors"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type info struct {
	symbol       *elf.Symbol
	section      *elf.Section
	symbolOffset uint64
}

const (
	program = "symbol-to-offset"
)

var (
	usage = fmt.Sprintf(`
Usage: %s EXECUTABLE SYMBOL

	EXECUTABLE	Path to the executable ELF file.
	SYMBOL		The name of the symbol in the executable file.
`, program)
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	info, err := symbolToOffset(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Symbol", "Symbol VA", "Symbol offset", "Section", "Section VA", "Section offset"})
	data := [][]string{{
		info.symbol.Name,
		fmt.Sprintf("%X", info.symbol.Value),
		fmt.Sprintf("%X", info.symbolOffset),
		info.section.Name,
		fmt.Sprintf("%X", info.section.Addr),
		fmt.Sprintf("%X", info.section.Offset),
	}}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// symbolToOffset attempts to resolve a symbol name to an offset in the binary.
func symbolToOffset(path, symbol string) (*info, error) {
	f, err := elf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open elf file to resolve symbol offset: %w", err)
	}

	regularSymbols, regularSymbolsErr := f.Symbols()
	dynamicSymbols, dynamicSymbolsErr := f.DynamicSymbols()

	if regularSymbolsErr != nil && dynamicSymbolsErr != nil {
		return nil, fmt.Errorf("could not open regular or dynamic symbol sections to resolve symbol offset: %w %s", regularSymbolsErr, dynamicSymbolsErr)
	}

	syms := append(regularSymbols, dynamicSymbols...)

	sectionsToSearchForSymbol := []*elf.Section{}

	for i := range f.Sections {
		if f.Sections[i].Flags == elf.SHF_ALLOC+elf.SHF_EXECINSTR {
			sectionsToSearchForSymbol = append(sectionsToSearchForSymbol, f.Sections[i])
		}
	}

	var executableSection *elf.Section

	for j := range syms {
		if syms[j].Name == symbol {
			for m := range sectionsToSearchForSymbol {
				if syms[j].Value > sectionsToSearchForSymbol[m].Addr &&
					syms[j].Value < sectionsToSearchForSymbol[m].Addr+sectionsToSearchForSymbol[m].Size {
					executableSection = sectionsToSearchForSymbol[m]
				}
			}

			if executableSection == nil {
				return nil, errors.New("could not find symbol in executable sections of binary")
			}

			return &info{
				symbol:       &syms[j],
				symbolOffset: syms[j].Value - executableSection.Addr + executableSection.Offset,
				section:      executableSection,
			}, nil
		}
	}

	return nil, fmt.Errorf("symbol %s not found in %s", symbol, path)
}
