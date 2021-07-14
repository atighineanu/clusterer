#### cookie:
The clusterer tool also has a non-related (to public cloud) feature "createvm" that creates virtual machines on a remote host (using ssh)
from the qcow2 images and xml vm description from http://rock229.qa.prv.suse.net (these XML files correspond to sle 15.1 - 15.3 OS versions
and are used for testing Xen virtualhosts as SUMA clients; they are already set up to boot from Xen vmlinuz image and XMLs are modified to 
seem like a bare-metal machine; more details inside the )