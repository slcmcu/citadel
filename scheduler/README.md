### Scheduler


### Events

* containers.start -  when a container starts on the cluster
* containers.stop - when a container is killed, dies, or stops

* slaves.joining - when a slave joins the cluster
* slaves.leaving - when a slave leaves the cluster
* slaves.pull - add a new image to the slaves

* execute.<slave_uuid> - execute a container on a slave
