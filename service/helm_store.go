package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"io"
	"k8sManagerApi/config"
	"k8sManagerApi/dao"
	"k8sManagerApi/model"
	"k8sManagerApi/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var HelmStore helmStore

type helmStore struct{}

// 定义列表返回的内容
type releaseElement struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Revision     string `json:"revision"`
	Updated      string `json:"updated"`
	Status       string `json:"status"`
	Chart        string `json:"chart"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`

	Notes string `json:"notes,omitempty"`
}

// releaseElements 定义返回的内容
type releaseElements struct {
	Items []*releaseElement `json:"items"`
	Total int               `json:"total"`
}

// ListRelease release列表，没有使用page和limit，前端实现
func (h *helmStore) ListRelease(actionConfig *action.Configuration, filterName string) (*releaseElements, error) {
	// new一个列表的Client
	client := action.NewList(actionConfig)
	client.Filter = filterName
	// 显示所有数据
	client.All = true
	//client.Limit = limit
	//client.Offset = offset
	client.TimeFormat = "2006-01-02 15:04:05"
	// 是否已经部署
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Release失败, %v", err.Error()))
		return nil, errors.New("获取Release失败," + err.Error())
	}
	total := len(results)
	elements := make([]*releaseElement, 0, len(results))
	for _, r := range results {
		elements = append(elements, constructReleaseElement(r, false))
	}
	releaseElementsRest := &releaseElements{
		Items: elements,
		Total: total,
	}
	return releaseElementsRest, nil
}

// DetailRelease release详情
func (h *helmStore) DetailRelease(actionConfig *action.Configuration, release string) (*release.Release, error) {
	client := action.NewGet(actionConfig)
	data, err := client.Run(release)
	if err != nil {
		zap.L().Error(fmt.Sprintf("获取Release详情失败, %v", err.Error()))
		return nil, errors.New("获取Release详情失败," + err.Error())
	}
	return data, nil
}

// InstallRelease 安装release, release: release的名字， chart：chart文件所在的路径, 上传
func (h *helmStore) InstallRelease(actionConfig *action.Configuration, cluster, release, chart, namespace string) error {
	client := action.NewInstall(actionConfig)
	client.ReleaseName = release
	// 这里的namespace没啥用，主要安装在哪个namespace还是要看actionConfig初始化的namespace
	client.Namespace = namespace
	splitChart := strings.Split(chart, ".")
	if splitChart[len(splitChart)-1] == "tgz" && strings.Contains(chart, ":") {
		chart = config.Conf.UploadPath + chart
	}
	// 加载chart文件，并基于文件内容生成k8的资源
	chart = filepath.Join(config.Conf.UploadPath, chart)
	chartRequested, err := loader.Load(chart)
	if err != nil {
		zap.L().Error(fmt.Sprintf("加载Chart文件, %v", err.Error()))
		return errors.New("加载Chart文件," + err.Error())
	}
	vals := make(map[string]interface{}, 0)
	_, err = client.Run(chartRequested, vals)
	if err != nil {
		zap.L().Error(fmt.Sprintf("安装Release文件, %v", err.Error()))
		return errors.New("安装Release文件," + err.Error())
	}
	return nil
}

// UninstallRelease 卸载Release
func (h *helmStore) UninstallRelease(actionConfig *action.Configuration, release, namespace string) error {
	client := action.NewUninstall(actionConfig)
	_, err := client.Run(release)
	if err != nil {
		zap.L().Error(fmt.Sprintf("卸载Release失败, %v", err.Error()))
		return errors.New("卸载Release失败," + err.Error())
	}
	return nil
}

// UploadChartFile Chart文件上传
func (h *helmStore) UploadChartFile(file multipart.File, header *multipart.FileHeader) error {
	filename := header.Filename
	t := strings.Split(filename, ".")
	if t[len(t)-1] != "tgz" {
		zap.L().Warn("Chart文件必须以.tgz结尾")
		return errors.New(fmt.Sprintf("Chart文件必须以.tgz结尾"))
	}

	filePath := config.Conf.UploadPath + filename
	_, err := os.Stat(filePath)
	if os.IsExist(err) {
		zap.L().Warn("Chart文件已存在")
		return errors.New(fmt.Sprintf("Chart文件已存在"))
	}
	out, err := os.Create(filePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建Chart文件失败, %v", err.Error()))
		return errors.New(fmt.Sprintf("创建Chart文件失败, %v", err.Error()))
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		zap.L().Error(fmt.Sprintf("创建Chart文件失败2, %v", err.Error()))
		return errors.New(fmt.Sprintf("创建Chart文件失败2, %v", err.Error()))
	}
	return nil
}

// DeleteChartFile Chart文件删除
func (h *helmStore) DeleteChartFile(chart string) error {
	//filePath := config.UploadPath + "/" + chart
	filePath := config.Conf.UploadPath + chart
	_, err := os.Stat(filePath)
	if err != nil || os.IsNotExist(err) {
		zap.L().Warn("Chart文件不存在")
		return errors.New(fmt.Sprintf("Chart文件不存在"))
	}
	err = os.Remove(filePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("删除Chart文件失败, %v", err.Error()))
		return errors.New(fmt.Sprintf("删除Chart文件失败, %v", err.Error()))
	}
	return nil
}

// 配置分发，现在没有使用
// scpChartToServer 拷贝chart到服务器，这里都是上传到临时目录/tmp/下面，并返回上传后的 chart在服务器上的路径
func (h *helmStore) scpChartToServer(cluster, chart string) (err error) {
	// 根据传入的cluster去数据库中查询对应的集群IP地址。
	ret, _ := dao.Node.GetNodeList(cluster)
	for _, v := range ret {
		// 循环集群中的每个node节点，将chart包传输到node节点的临时目录/tmp下面
		sftp := utils.Cli{
			Username: "root",
			Password: "ives",
			Host:     v.IP,
			Port:     22,
			Client:   nil,
		}
		// 进行连接
		err := sftp.Connection()
		if err != nil {
			zap.L().Error("连接集群的node节点失败")
		}
		// 文件拷贝
		srcChartPath := filepath.Join(config.Conf.UploadPath, chart)
		destChartPath := filepath.Join("/tmp", chart)
		err = sftp.ScpToServer(srcChartPath, destChartPath)
		if err != nil {
			zap.L().Error("向node节点传递chart包失败")
		}
	}
	return nil
}

// GetCharts 获取Chart列表
func (h *helmStore) GetCharts(name string, page, limit int) (*dao.Charts, error) {
	return dao.Chart.GetList(name, page, limit)
}

// AddChart Chart新增
func (h *helmStore) AddChart(chart *model.Chart) error {
	_, has, err := dao.Chart.HasChart(chart.Name)
	if err != nil {
		return err
	}
	if has {
		return errors.New(fmt.Sprintf("该数据已存在，请重新添加"))
	}
	if err := dao.Chart.Add(chart); err != nil {
		return err
	}
	return nil
}

// UpdateChart Chart更新
func (h *helmStore) UpdateChart(chart *model.Chart) error {
	oldChart, _, err := dao.Chart.HasChart(chart.Name)
	if err != nil {
		return err
	}
	fmt.Println(chart.FileName, oldChart.FileName)
	// 这里先判断，如果文件为空和不和之前的相同，则删除之前的文件
	if chart.FileName != "" && chart.FileName != oldChart.FileName {
		err = h.DeleteChartFile(oldChart.FileName)
		if err != nil {
			return err
		}
	}
	return dao.Chart.Update(chart)
}

// DeleteChart Chart删除
func (h *helmStore) DeleteChart(chart *model.Chart) error {
	// 删除文件
	err := h.DeleteChartFile(chart.FileName)
	if err != nil {
		return err
	}
	// 删除数据
	return dao.Chart.Delete(chart.ID)
}

// constructReleaseElement release内容过滤
func constructReleaseElement(r *release.Release, showStatus bool) *releaseElement {
	element := &releaseElement{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Revision:     strconv.Itoa(r.Version),
		Status:       r.Info.Status.String(),
		Chart:        r.Chart.Metadata.Name,
		ChartVersion: r.Chart.Metadata.Version,
		AppVersion:   r.Chart.Metadata.AppVersion,
	}
	if showStatus {
		element.Notes = r.Info.Notes
	}
	// 输出判断
	t := "-"
	if tspb := r.Info.LastDeployed; !tspb.IsZero() {
		t = tspb.String()
	}
	element.Updated = t
	return element
}