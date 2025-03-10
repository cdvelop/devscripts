#!/bin/bash
# VPS security setup script for Debian-based Linux servers
# Usage: sudo ./vpssetupsecurity.sh <username> <new_ssh_port>

# Import helper functions
source functions.sh

# Check if script is run as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root"
        exit 1
    fi
}

# Configure SSH security settings
configure_ssh_security() {
    local username=$1
    
    warning "Configuring SSH security settings..."
    
    # Disabling root user login
    execute "sed -i 's/^#PermitRootLogin prohibit-password$/PermitRootLogin no/' /etc/ssh/sshd_config" \
            "Failed to disable root login" \
            "Root login disabled successfully"
    
    # Set max auth tries
    execute "sed -i 's/^#MaxAuthTries 6$/MaxAuthTries 3/' /etc/ssh/sshd_config" \
            "Failed to set MaxAuthTries" \
            "MaxAuthTries set to 3 successfully"
    
    # Set max sessions
    execute "sed -i 's/^#MaxSessions 10$/MaxSessions 3/' /etc/ssh/sshd_config" \
            "Failed to set MaxSessions" \
            "MaxSessions set to 3 successfully"
    
    # Set login grace time
    execute "sed -i 's/^#LoginGraceTime 2m$/LoginGraceTime 30/' /etc/ssh/sshd_config" \
            "Failed to set LoginGraceTime" \
            "LoginGraceTime set to 30 seconds successfully"
    
    # Allow only specified user
    execute "echo \"AllowUsers $username\" >> /etc/ssh/sshd_config" \
            "Failed to set AllowUsers" \
            "AllowUsers set to $username successfully"
    
    return 0
}

# Change SSH port
change_ssh_port() {
    local new_port=$1
    
    warning "Changing SSH port to $new_port..."
    
    execute "grep ssh /etc/services" \
            "Failed to check SSH services" \
            "SSH services checked successfully"
    
    execute "sed -i \"s/^#Port 22\$/Port $new_port/\" /etc/ssh/sshd_config" \
            "Failed to change SSH port" \
            "SSH port changed to $new_port successfully"
    
    execute "systemctl restart sshd" \
            "Failed to restart SSH service" \
            "SSH service restarted successfully"
    
    return 0
}

# Setup firewall
setup_firewall() {
    local new_port=$1
    
    warning "Setting up firewall..."
    
    execute "apt install firewalld -y" \
            "Failed to install firewalld" \
            "Firewalld installed successfully"
    
    execute "systemctl enable firewalld" \
            "Failed to enable firewalld service" \
            "Firewalld service enabled successfully"
    
    execute "firewall-cmd --permanent --zone=public --add-port=\"$new_port\"/tcp" \
            "Failed to add SSH port to firewall" \
            "SSH port $new_port added to firewall"
    
    execute "firewall-cmd --permanent --zone=public --add-port=443/tcp" \
            "Failed to add HTTPS port to firewall" \
            "HTTPS port 443 added to firewall"
    
    execute "firewall-cmd --reload" \
            "Failed to reload firewall configuration" \
            "Firewall configuration reloaded successfully"
    
    return 0
}

# Verify service status
verify_services() {
    local new_port=$1
    
    warning "Verifying services status..."
    
    execute "systemctl status sshd --no-pager" \
            "Failed to check SSH status" \
            "SSH status checked successfully"
    
    execute "firewall-cmd --state" \
            "Failed to check firewall state" \
            "Firewall state checked successfully"
    
    execute "systemctl status firewalld --no-pager" \
            "Failed to check firewalld service status" \
            "Firewalld service status checked successfully"
    
    execute "ss -an | grep \"$new_port\"" \
            "Failed to verify SSH daemon is listening on port $new_port" \
            "SSH daemon is listening on port $new_port"
    
    return 0
}

# Main execution function
main() {
    local username=$1
    local new_port=$2
    
    check_root
    
    warning "==== Starting VPS Security Configuration ===="
    
    # Configure SSH security
    configure_ssh_security "$username"
    
    # Change SSH port
    change_ssh_port "$new_port"
    
    # Setup firewall
    setup_firewall "$new_port"
    
    # Verify services are running correctly
    verify_services "$new_port"
    
    warning "==== VPS Security Configuration Completed ===="
    success "IMPORTANT: Connect using port $new_port for future SSH sessions"
    
    # Print all success messages
    successMessages
}

# Validate parameters
if [ -z "$1" ] || [ -z "$2" ]; then
    error "Missing required parameters"
    echo "Usage: sudo $0 <username> <new_ssh_port>"
    echo "Example: sudo $0 admin 2222"
    exit 1
fi

# Run the main function
main "$1" "$2"
