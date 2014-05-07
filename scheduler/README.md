### Scheduler

The scheduler is a very simple interface that takes `profiles` and looks at the current
resources used by the host and the resources requested for a container to determine 
if a host can run a container or not.

Profiles are written in lua and take two parameters.  `current` and `requested`.  The `current` 
parameter provides the current information of the host.  The `requested` parameter provides the
resources requested to run a container.  The profile must defind one method named `Accept(current, requested)` that returns true or false if the host can run the container.

A sample profile looks like this:

```lua
-- yes profile
function Accept(current, requested)
    return true
end
```

The requested and current types have the following data:
* Cpus - the number of cpus required 
* CpuProfile - high,medium,low defining the cpu usage
* Memory - the amount of memory required
* Disk - the about of disk space required
* Image - the image of the container to be run
* Context - generic map of data
