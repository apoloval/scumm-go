package scumm

import "github.com/apoloval/scumm-go/vm"

func Run(indexPath string) error {
	rm, err := FromIndexFile(indexPath)
	if err != nil {
		return err
	}

	return vm.NewEngine(rm).Run()
}
