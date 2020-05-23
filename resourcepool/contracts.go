package resourcepool

type bootable interface {
	boot(resourcePool *ResourcePool) error
}
