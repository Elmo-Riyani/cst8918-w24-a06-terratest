package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "1cf22e62-8dc8-4849-b953-9b546c070c52"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "riya0001",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
	// 2️⃣ Confirm NIC Exists & Is Connected to VM
	nic, err := azure.GetNetworkInterfaceE(nicName, resourceGroupName, subscriptionID)
	assert.NoError(t, err, "Failed to get network interface")
	assert.NotNil(t, nic.VirtualMachine, "NIC is not connected to any VM")
	assert.Contains(t, *nic.VirtualMachine.ID, vmName, "NIC is not connected to the expected VM")

	// 3️⃣ Confirm the VM is Running the Correct Ubuntu Version
	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)
	expectedSku := "22_04-lts-gen2"
	actualSku := *vm.StorageProfile.ImageReference.Sku
	assert.Equal(t, expectedSku, actualSku, "VM is not running expected Ubuntu SKU version")
}
