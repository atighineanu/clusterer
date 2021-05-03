#### *CLUSTERER*  - a lightweight tool that uses libvirt - related (os) commands to manage VMs locally or on a remote VirtHost.

For cloning & creating a "fake" bare-metal Xen VirtHost (Which will be a SUMA Client)... on your remote Virtualization Host - run
`clusterer fakexenbm --rmtip <val> --distro "<val>"`

*fakexenbm* - from fake xen bare-metal 

*--rmtip* - from "remote IP" - the ip of your remote VM Host/ could be one of the "rock" machines

*--distro* - from (sles) "distribution" sles distro with possible values '15.1', '15.2', '15.3'

for example you can run: `clusterer fakexenbm --rmtip 10.84.154.100 --distro 15.3`
this will copy a compressed qcow2 image, and an xml-vm-profile from  10.84.149.229/ on your remote machine, will create the vm and start it.
