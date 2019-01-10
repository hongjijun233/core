package docker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"

	corecluster "github.com/projecteru2/core/cluster"
	coretypes "github.com/projecteru2/core/types"
	log "github.com/sirupsen/logrus"

	"github.com/docker/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/registry"
)

type fuckDockerStream struct {
	conn net.Conn
	buf  io.Reader
}

func (f fuckDockerStream) Read(p []byte) (n int, err error) {
	return f.Read(p)
}

func (f fuckDockerStream) Close() error {
	return f.conn.Close()
}

func doFuckDockerStream(stream io.ReadCloser) io.Reader {
	outr, outw := io.Pipe()

	go func() {
		defer stream.Close()
		_, err := stdcopy.StdCopy(outw, outw, stream)
		outw.CloseWithError(err)
	}()

	return outr
}

// FuckDockerStream will copy docker stream to stdout and err
func FuckDockerStream(stream dockertypes.HijackedResponse) io.ReadCloser {
	outr := doFuckDockerStream(ioutil.NopCloser(stream.Reader))
	return fuckDockerStream{stream.Conn, outr}
}

// make mount paths
// 使用volumes, 参数格式跟docker一样
// volumes:
//     - "/foo-data:$SOMEENV/foodata:rw"
func makeMountPaths(opts *coretypes.DeployOptions) ([]string, map[string]struct{}) {
	binds := []string{}
	volumes := make(map[string]struct{})

	var expandENV = func(env string) string {
		envMap := map[string]string{}
		for _, env := range opts.Env {
			parts := strings.Split(env, "=")
			envMap[parts[0]] = parts[1]
		}
		return envMap[env]
	}

	for _, path := range opts.Volumes {
		expanded := os.Expand(path, expandENV)
		parts := strings.Split(expanded, ":")
		if len(parts) == 2 {
			binds = append(binds, fmt.Sprintf("%s:%s:rw", parts[0], parts[1]))
			volumes[parts[1]] = struct{}{}
		} else if len(parts) == 3 {
			binds = append(binds, fmt.Sprintf("%s:%s:%s", parts[0], parts[1], parts[2]))
			volumes[parts[1]] = struct{}{}
		}
	}

	return binds, volumes
}

func makeResourceSetting(cpu float64, memory int64, cpuMap map[string]int, softlimit bool) dockercontainer.Resources {
	resource := dockercontainer.Resources{}
	if cpu > 0 {
		resource.CPUPeriod = corecluster.CPUPeriodBase
		resource.CPUQuota = int64(cpu * float64(corecluster.CPUPeriodBase))
	}
	if cpuMap != nil && len(cpuMap) > 0 {
		cpuIDs := []string{}
		for cpuID := range cpuMap {
			cpuIDs = append(cpuIDs, cpuID)
		}
		resource.CpusetCpus = strings.Join(cpuIDs, ",")
	}
	if softlimit {
		resource.MemoryReservation = memory
	} else {
		resource.Memory = memory
		resource.MemorySwap = memory
		if memory/2 > minMemory {
			resource.MemoryReservation = memory / 2
		}
	}
	return resource
}

// 只要一个image的前面, tag不要
func normalizeImage(image string) string {
	if strings.Contains(image, ":") {
		t := strings.Split(image, ":")
		return t[0]
	}
	return image
}

// image begin
// MakeAuthConfigFromRemote Calculate encoded AuthConfig from registry and eru-core config
// See https://github.com/docker/cli/blob/16cccc30f95c8163f0749eba5a2e80b807041342/cli/command/registry.go#L67
func makeEncodedAuthConfigFromRemote(authConfigs map[string]coretypes.AuthConfig, remote string) (string, error) {
	ref, err := reference.ParseNormalizedNamed(remote)
	if err != nil {
		return "", err
	}

	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return "", err
	}

	serverAddress := repoInfo.Index.Name
	if authConfig, exists := authConfigs[serverAddress]; exists {
		if encodedAuth, err := encodeAuthToBase64(authConfig); err == nil {
			return encodedAuth, nil
		}
		return "", err
	}
	return "dummy", nil
}

// EncodeAuthToBase64 serializes the auth configuration as JSON base64 payload
// See https://github.com/docker/cli/blob/master/cli/command/registry.go#L41
func encodeAuthToBase64(authConfig coretypes.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

func getIP(daemonHost string) string {
	u, err := url.Parse(daemonHost)
	if err != nil {
		log.Errorf("[getIP] GetIP %s failed %v", daemonHost, err)
		return ""
	}
	return u.Hostname()
}