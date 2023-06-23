# Create a build environment for Docker plugin
require "yaml"
settings = YAML.load_file "settings.yaml"

IP_SECTIONS = settings["network"]["control_ip"].match(/^([0-9.]+\.)([^.]+)$/)
# First 3 octets including the trailing dot:
IP_NW = IP_SECTIONS.captures[0]
# Last octet excluding all dots:
IP_START = Integer(IP_SECTIONS.captures[1])

Vagrant.configure("2") do |config|
  config.vm.box = settings["software"]["box"]

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  config.vm.box_check_update = true
  config.vm.network "private_network", ip: settings["network"]["control_ip"]

  config.vm.define "devnerdctld" do |devnerdctld|
    devnerdctld.vm.hostname = "devnerdctld"
    
    # Configure resource
    devnerdctld.vm.provider "virtualbox" do |vb|
      vb.cpus = settings["nodes"]["workers"]["cpu"]
      vb.memory = settings["nodes"]["workers"]["memory"]
    end
  end

  # Enable provisioning with a shell script. Additional provisioners such as
  # Ansible, Chef, Docker, Puppet and Salt are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", 
    path: "scripts/init.sh"
    # path: "scripts/init_podman.sh"
end
