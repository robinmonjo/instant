Vagrant::Config.run do |config|

  config.vm.box = "precise64"
  config.vm.box_url = "http://files.vagrantup.com/precise64.box"

  config.vm.forward_port 5000, 5000

  config.vm.provision :shell, :inline => <<EOF
/vagrant/setup-host.sh
EOF

end

Vagrant.configure("2") do |config|
  config.vm.network :public_network
end