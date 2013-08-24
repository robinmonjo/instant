Vagrant::Config.run do |config|

  config.vm.box = "precise64"
  config.vm.box_url = "http://files.vagrantup.com/precise64.box"

  config.vm.provision :shell, :inline => <<EOF
/go/src/github.com/rmonjo/instant/setup-host.sh
EOF

end

Vagrant.configure("2") do |config|
  config.vm.network :public_network
  config.vm.synced_folder "../../../../", "/go"
end