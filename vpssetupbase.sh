#!/bin/bash
# Base VPS setup for Debian-based Linux servers
# Usage: sudo ./vpssetupbase.sh <username> <ssh_key>

# Import helper functions
source functions.sh

# Check if script is run as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root" 
        exit 1
    fi
}

# Configure SSH directory and keys
setup_ssh() {
    local username=$1
    local ssh_key=$2
    
    warning "Setting up SSH directory and keys..."
    
    execute "mkdir -p ~/.ssh" \
            "Failed to create SSH directory" \
            "SSH directory created successfully"
    
    execute "chmod 700 ~/.ssh" \
            "Failed to set permissions on SSH directory" \
            "SSH directory permissions set correctly"
    
    execute "echo \"$ssh_key\" > ~/.ssh/authorized_keys" \
            "Failed to create authorized_keys file" \
            "Authorized keys file created successfully"
    
    execute "chmod 600 ~/.ssh/authorized_keys" \
            "Failed to set permissions for authorized_keys" \
            "Authorized keys permissions set correctly"
    
    return 0
}

# Set system timezone
set_timezone() {
    local timezone=$1
    
    warning "Setting timezone to $timezone..."
    
    execute "timedatectl set-timezone \"$timezone\"" \
            "Failed to set timezone" \
            "Timezone set to $timezone successfully"
    
    return 0
}

# Remove ubuntu default user if it exists
remove_ubuntu_user() {
    warning "Checking for ubuntu default user..."
    
    if id "ubuntu" &>/dev/null; then
        warning "Removing ubuntu user with directory"
        execute "userdel -r -f ubuntu" \
                "Failed to remove ubuntu user" \
                "User ubuntu removed successfully"
        
        execute "deluser --system --remove-home ubuntu" \
                "Failed to delete ubuntu system user" \
                "System user ubuntu removed successfully" \
                "continue_on_error"
    else
        success "Ubuntu user not found, skipping removal."
    fi
    
    return 0
}

# Check user and groups
check_user() {
    local username=$1
    
    warning "Checking user groups for $username..."
    
    execute "groups \"$username\"" \
            "Failed to check groups for $username" \
            "User groups for $username retrieved successfully"
    
    return 0
}

# Main execution function
main() {
    local username=$1
    local ssh_key=$2
    local timezone=${3:-"America/Santiago"}
    
    check_root
    
    warning "==== Starting Basic VPS Setup ===="
    
    # Setup SSH
    setup_ssh "$username" "$ssh_key"
    
    # Set timezone
    set_timezone "$timezone"
    
    # Remove default ubuntu user if present
    remove_ubuntu_user
    
    # Check user and connected users
    check_user "$username"
    
    warning "Checking connected users..."
    execute "w" \
            "Failed to check connected users" \
            "Connected users retrieved successfully"
    
    success "==== Basic VPS Setup Completed Successfully ===="
    
    # Print all success messages
    successMessages
}

# Validate parameters
if [ -z "$1" ] || [ -z "$2" ]; then
    error "Missing required parameters" 
    echo "Usage: sudo $0 <username> <ssh_key> [timezone]"
    echo "Example: sudo $0 admin 'ssh-rsa AAAAB3NzaC1...' 'America/Santiago'"
    exit 1
fi

# Run the main function
main "$1" "$2" "$3"
