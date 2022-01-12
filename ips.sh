#!/bin/bash
a=`sudo virsh list --all | grep running | cut -d " " -f 5,6 | grep -i -v seed`; 

case ${1} in
	"")
	for i in ${a}; do echo "${i}: " ; sudo virsh domifaddr ${i} | grep ipv4 | cut -d " " -f21,20,22 | cut -d"/" -f1 ;  done
	;;

	"q")
	for i in ${a}; do sudo virsh domifaddr ${i} | grep ipv4 | cut -d " " -f21,20,22 | cut -d"/" -f1 ;  done
	;;

	"j")
	echo "{"
	declare -i count;
	declare -i tmp=`echo ${a} | wc -w`; 
	for i in ${a}
	do 
	{
	count=$((count+1)); name=`echo "${i}" | cut -d":" -f1 | sed s/default-//`; b=`sudo virsh domifaddr ${i} | grep ipv4 | cut -d " " -f21,20,22 | cut -d"/" -f1 | cut -d" " -f2 `; 
	if (( ${count} < ${tmp} )) 
	then 
		{
			echo "\"${name}\": \"${b}\","; 
		} else 
			{
				echo "\"${name}\": \"${b}\"";
			 } 
	fi  
	}
	done
	echo "}"
	;;
	"l")
	for i in ${a}
	do 
	{
	name=`echo "${i}" | cut -d":" -f1 | sed s/default-//`; b=`sudo virsh domifaddr ${i} | grep ipv4 | cut -d " " -f21,20,22 | cut -d"/" -f1 | cut -d" " -f2 `; 
	echo "${name}: ${b}"
	}
	done
	;;
	"i")
	echo '[all]' > ./ansible/inventory
	for i in ${a}
	do 
	{
	name=`echo "${i}" | cut -d":" -f1 | sed s/default-//`; b=`sudo virsh domifaddr ${i} | grep ipv4 | cut -d " " -f21,20,22 | cut -d"/" -f1 | cut -d" " -f2 `; 
	echo "${b}	hostname=${name}" >> ./ansible/inventory
	}
	done
	;;
esac
	

