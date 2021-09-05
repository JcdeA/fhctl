package choose

import (
	"fmt"

	"github.com/jcdea/aarch64-client-go"
	"github.com/jcdea/fhctl/check"
	"github.com/jcdea/fhctl/request"
	"github.com/jcdea/fhctl/types"
	"github.com/manifoldco/promptui"
)

func ChooseProxy(args []string) (vm aarch64.VM, err error) {
	// We already check if user is signed in at GetProjects
	var VMlist []aarch64.VM

	if len(args) != 0 && args[0] != "" {
		VMlist, err = withAlias(args)
		if err != nil {
			return chooseVMFromList(VMlist)
		}

		println(fmt.Sprintf("Alias not found for\"%v\"\n", args[0]))
	}

	projectResp, err := request.GetProjects()
	check.CheckErr(err, "Failed to retrieve project list")

	var projectNames []string
	for _, item := range projectResp.Projects {
		projectNames = append(projectNames, item.Name)
	}

	check.CheckErr(err, "")
	selectProject := promptui.Select{
		Label: "Select project",
		Items: projectNames,
	}
	selectedProjectIndex, _, err := selectProject.Run()
	check.CheckErr(err, "")

	selectedProject := projectResp.Projects[selectedProjectIndex]
	VMlist = selectedProject.VMs

	return chooseProxyFromList(VMlist)

}

// Search for alias, then ssh if alias is found.
// if not found: returns not found error
func withAlias(args []string) ([]aarch64.VM, error) {
	var vms []aarch64.VM

	project, err := types.SearchProjectAlias(args[0])
	if err != nil {
		return []aarch64.VM{}, err
	}

	println(fmt.Sprintf("Using alias %v=%v\n", args[0], project.Id))

	projectResp, err := request.GetProjects()
	check.CheckErr(err, "Failed to retrieve project list")

	for _, item := range projectResp.Projects {
		if item.Id == project.Id {
			vms = item.VMs

		}

	}
	return vms, nil
}

// choose the VM from list of VMs in project
func chooseProxyFromList(VMlist []aarch64.VM) (aarch64.VM, error) {
	var VMNames []string
	for _, vm := range VMlist {
		VMNames = append(VMNames, formatVMmeta(vm))
	}

	selectVM := promptui.Select{
		Label: "Select VM",
		Items: VMNames,
	}
	selectedVMIndex, _, err := selectVM.Run()
	check.CheckErr(err, "")
	return VMlist[selectedVMIndex], nil

}
