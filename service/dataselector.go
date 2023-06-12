package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

// dataSelector 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell 接口，用于各种资源List的类型转换，转换后可以使用dataSelector的排序、过滤、分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的结构体，过滤：Name，分页：Limit和Page
type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// 实现自定义结构的排序，需要重写Len、Swap、Less方法
// Len 用于获取数组的长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 用于数据比较大小之后的位置变更
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 用于排序，比较大小
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
	//return d.GenericDataList[i].GetCreation().Before(d.GenericDataList[j].GetCreation())
}

// Sort 重写以上三个方法，使用sort.Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// Filter 方法用于过滤数据，比较数据的Name属性，若包含，则返回，
func (d *dataSelector) Filter() *dataSelector {
	// 判断参数name是否为空，若为空，则返回所有数据
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	//	如果不为空，则按照入参Name进行过滤，然后返回
	filtered := []DataCell{}
	for _, value := range d.GenericDataList {
		// 定义是否匹配的标签变量，默认是匹配
		match := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			match = false
			continue
		}
		if match {
			filtered = append(filtered, value)
		}
	}
	d.GenericDataList = filtered
	return d
}

// Paginate 方法用于分页，获取分页数据, 根据Limit和Page的传惨，取一定范围内的数据，返回
func (d *dataSelector) Paginate() *dataSelector {
	// 根据Limit和Page的入参，定义快捷变量
	limit := d.DataSelect.Paginate.Limit
	page := d.DataSelect.Paginate.Page
	// 检验参数的合法性
	if limit <= 0 || page <= 0 {
		return d
	}
	// 定义取数据范围需要的startIndex 和 endIndex
	//举例：25个元素的数组，limit是10，page是3，startIndex是20，endIndex是30（实际上endIndex是 24）
	startIndex := limit * (page - 1)
	endIndex := limit * page
	// 处理endIndex
	if endIndex > len(d.GenericDataList) {
		endIndex = len(d.GenericDataList)
	}
	// 获取分页数据
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// podCell 定义podCell 重写GetCreation和GetName方法后，可以进行数据交换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

// deploymentCell 定义deploymentCell 重写GetCreation和GetName方法后，可以进行数据交换
type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

// daemonSetCell 定义DaemonSetCell 重写GetCreation和GetName方法后，可以进行数据交换
type daemonSetCell appsv1.DaemonSet

func (d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonSetCell) GetName() string {
	return d.Name
}

// statefulSetCell 定义StatefulSetCell 重写GetCreation和GetName方法后，可以进行数据交换
type statefulSetCell appsv1.StatefulSet

func (s statefulSetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulSetCell) GetName() string {
	return s.Name
}

// nodeCell 定义nodeCell 重写GetCreation和GetName方法后，可以进行数据交换
type nodeCell corev1.Node

func (n nodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nodeCell) GetName() string {
	return n.Name
}

// namespaceCell 定义namespaceCell 重写GetCreation和GetName方法后，可以进行数据交换
type namespaceCell corev1.Namespace

func (n namespaceCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n namespaceCell) GetName() string {
	return n.Name
}

// pvCell 定义pvCell 重写GetCreation和GetName方法后，可以进行数据交换
type pvCell corev1.PersistentVolume

func (p pvCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvCell) GetName() string {
	return p.Name
}

// serviceCell 定义serviceCell 重写GetCreation和GetName方法后，可以进行数据交换
type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}

// serviceCell 定义serviceCell 重写GetCreation和GetName方法后，可以进行数据交换
type ingressCell nwv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}

// configmapCell 定义configmapCell 重写GetCreation和GetName方法后，可以进行数据交换
type configmapCell corev1.ConfigMap

func (c configmapCell) GetCreation() time.Time {
	return c.CreationTimestamp.Time
}

func (c configmapCell) GetName() string {
	return c.Name
}

// secretCell 定义secretCell 重写GetCreation和GetName方法后，可以进行数据交换
type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s secretCell) GetName() string {
	return s.Name
}

// pvcCell 定义pvcCell 重写GetCreation和GetName方法后，可以进行数据交换
type pvcCell corev1.PersistentVolumeClaim

func (p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvcCell) GetName() string {
	return p.Name
}