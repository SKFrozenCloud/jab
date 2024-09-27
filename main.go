package main

import (
	"fmt"
	"time"
)

var (
	DBFile                       = "hashes.db"
	AESKey                       = ""
	CheckIntervalSeconds         = 60
	SensitiveFilesAndDirectories = []string{
		// System Configuration Files
		"/etc/passwd",
		"/etc/shadow",
		"/etc/group",
		"/etc/sudoers",
		"/etc/hosts",
		"/etc/hostname",
		"/etc/ssh/sshd_config",
		"/etc/ssh/ssh_config",
		"/etc/fstab",
		"/etc/sysctl.conf",
		"/etc/crontab",
		"/etc/cron.*/*",
		"/etc/resolv.conf",
		"/etc/nsswitch.conf",
		"/etc/pam.d/*",
		"/etc/security/*",
		// Kernel and Boot Configuration
		"/boot/grub/grub.cfg",
		"/boot/vmlinuz-*",
		"/boot/initrd.img-*",
		// Network Configuration
		"/etc/network/interfaces",
		"/etc/netplan/*",
		"/etc/sysconfig/network-scripts/*",
		"/etc/iptables/*",
		"/etc/firewalld/*",
		// SSH and Authorized Keys
		"/root/.ssh/authorized_keys",
		"/home/*/.ssh/authorized_keys",
		"/home/*/.bash_history",
		"/home/*/.bashrc",
		"/home/*/.profile",
		// Application and Service Configuration
		"/etc/apache2/*",
		"/etc/httpd/*",
		"/etc/nginx/*",
		"/etc/mysql/my.cnf",
		"/etc/my.cnf",
		"/etc/postgresql/*",
		"/etc/redis/redis.conf",
		"/etc/samba/smb.conf",
		"/etc/mail/*",
		// Sensitive User Files
		"/root/.bashrc",
		"/root/.profile",
		"/root/.history",
		"/home/*/.config/*",
		// Other Sensitive Files and Directories
		"/etc/hosts.allow",
		"/etc/hosts.deny",
		"/etc/gshadow",
		"/etc/ld.so.conf",
		"/proc/sys/net/*",
		"/var/spool/cron/crontabs/*",
	}
)

func main() {
	fmt.Println("Starting monitoring file integrity...")

	for {
		// Running
		fmt.Println("Checking integrity...")
		integrityChanges, err := CheckIntegrity(DBFile, SensitiveFilesAndDirectories)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		// Result
		fmt.Println("Files added:")
		for i, v := range integrityChanges.Added {
			fmt.Printf("Added file number %v; Path: %v\n", i, v.FilePath)
		}
		fmt.Println("Files modified:")
		for i, v := range integrityChanges.Modified {
			fmt.Printf("Modified file number %v; Path: %v\n", i, v.FilePath)
		}
		fmt.Println("Files removed:")
		for i, v := range integrityChanges.Removed {
			fmt.Printf("Removed file number %v; Path: %v\n", i, v.FilePath)
		}

		time.Sleep(time.Duration(CheckIntervalSeconds) * time.Second)
	}
}
