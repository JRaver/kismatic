package integration_tests

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func runViaSSH(cmds []string, hosts []NodeDeets, sshKey string, period time.Duration) error {
	timeout := time.After(period)
	bail := make(chan struct{})
	cmdSuccess := make(chan bool)
	// Create a goroutine per host. Each goroutine runs the commands serially on the host
	// until one of these is true:
	// a) all commands were executed successfully,
	// b) an error occurred when running a command,
	// c) the goroutine got a signal to bail
	for _, host := range hosts {
		go func(node NodeDeets) {
			for _, cmd := range cmds {
				res, err := executeCmd(cmd, node.PublicIP, node.SSHUser, sshKey)
				fmt.Println(res)
				select {
				case cmdSuccess <- err == nil:
					if err != nil {
						return
					}
				case <-bail:
					return
				}
			}
		}(host)
	}

	// The bail channel is closed if we encounter an error, or if the timeout is reached.
	// This will signal all goroutines to return.
	defer close(bail)

	// At most, we will get a total of hosts * cmds status messages in the channel.
	for i := 0; i < len(hosts)*len(cmds); i++ {
		select {
		case ok := <-cmdSuccess:
			if !ok {
				return fmt.Errorf("error running command on node")
			}
		case <-timeout:
			return fmt.Errorf("timed out running commands on nodes")
		}
	}
	return nil
}

func executeCmd(cmd, hostname, user, sshKey string) (string, error) {
	sshCmd := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "-t", "-t", "-i", sshKey, user+"@"+hostname, cmd)
	sshCmd.Stdin = os.Stdin
	sshOut, sshErr := sshCmd.CombinedOutput()
	return hostname + ": " + string(sshOut), sshErr
}

func copyFileToRemote(file string, destFile string, node NodeDeets, sshKey string, period time.Duration) error {
	timeout := time.After(period)
	success := make(chan bool)
	go func() {
		out, err := scpFile(file, destFile, node.SSHUser, node.PublicIP, sshKey)
		fmt.Println(out)
		success <- err == nil
	}()
	select {
	case ok := <-success:
		if !ok {
			return errors.New("failed to copy file to node")
		}
	case <-timeout:
		return errors.New("timed out copying file to node")
	}
	return nil
}

func scpFile(filePath string, destFilePath string, user, hostname, sshKey string) (string, error) {
	ver := exec.Command("scp", "-o", "StrictHostKeyChecking no", "-i", sshKey, filePath, user+"@"+hostname+":"+destFilePath)
	out, err := ver.CombinedOutput()
	return string(out), err
}

// WaitUntilSSHOpen waits up to the given timeout for a successful SSH connection to
// the given node. If the connection is open, returns true. If the timeout is reached, returns false.
func WaitUntilSSHOpen(publicIP, sshUser, sshKey string, timeout time.Duration) bool {
	tout := time.After(timeout)
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-tout:
			return false
		case <-tick:
			cmd := exec.Command("ssh")
			cmd.Args = append(cmd.Args, "-i", sshKey)
			cmd.Args = append(cmd.Args, "-o", "ConnectTimeout=5")
			cmd.Args = append(cmd.Args, "-o", "BatchMode=yes")
			cmd.Args = append(cmd.Args, "-o", "StrictHostKeyChecking=no")
			cmd.Args = append(cmd.Args, fmt.Sprintf("%s@%s", sshUser, publicIP), "exit") // just call exit if we are able to connect
			if err := cmd.Run(); err == nil {
				// command succeeded
				fmt.Println()
				return true
			}
			fmt.Printf("?")
		}
	}
}
