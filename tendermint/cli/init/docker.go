package init

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func getNodeDockerComposeConfig(nodeIndex int, output io.Writer) {
	t, err := template.New("docker-compose-config").Parse(`
  abci-{{.NodeIndex}}:
    <<: *abci-node
    container_name: likechain_abci-{{.NodeIndex}}
  tendermint-{{.NodeIndex}}:
    <<: *tendermint-node
    container_name: likechain_tendermint-{{.NodeIndex}}
    hostname: tendermint-{{.NodeIndex}}
    depends_on:
      - abci-{{.NodeIndex}}
    volumes:
      - ./tendermint/nodes/{{.NodeIndex}}:/tendermint
    ports:
      - {{.P2PPort}}:26656
      - {{.RPCPort}}:26657
    command:
      - --proxy_app=tcp://abci-{{.NodeIndex}}:26658`)
	if err != nil {
		panic(err)
	}
	type templateParams struct {
		NodeIndex int
		P2PPort   int
		RPCPort   int
	}
	t.Execute(output, templateParams{
		NodeIndex: nodeIndex,
		P2PPort:   26658 + nodeIndex*2,
		RPCPort:   26659 + nodeIndex*2,
	})
}

func genDockerComposeFile(sampleInputPath, outputPath string, pubInfos []publicInfo) {
	sample, err := ioutil.ReadFile(sampleInputPath)
	if err != nil {
		panic(err)
	}
	buf := strings.Builder{}
	buf.Write(sample)
	for nodeIndex := 2; nodeIndex <= len(pubInfos); nodeIndex++ {
		getNodeDockerComposeConfig(nodeIndex, &buf)
	}
	outputString := buf.String()
	buf.Reset()
	buf.WriteString(fmt.Sprintf("%s@tendermint-1:26656", string(pubInfos[0].NodeID)))
	for i, pubInfo := range pubInfos[1:] {
		buf.WriteString(fmt.Sprintf(",%s@tendermint-%d:26656", string(pubInfo.NodeID), i+2))
	}
	peersString := buf.String()
	outputString = strings.Replace(outputString, "$PERSISTENT_PEERS", peersString, 1)
	err = ioutil.WriteFile(outputPath, []byte(outputString), 0644)
	if err != nil {
		panic(err)
	}
}

func genDockerComposeFiles(dockerDir string, pubInfos []publicInfo) {
	dockerDevSamplePath := fmt.Sprintf("%s/docker-compose.sample.yml", dockerDir)
	_, err := os.Stat(dockerDevSamplePath)
	if os.IsNotExist(err) {
		panic("Cannot find docker-compose.sample.yml")
	}
	dockerDevOutputPath := fmt.Sprintf("%s/docker-compose.yml", dockerDir)
	genDockerComposeFile(dockerDevSamplePath, dockerDevOutputPath, pubInfos)

	dockerProdSamplePath := fmt.Sprintf("%s/docker-compose.production.sample.yml", dockerDir)
	_, err = os.Stat(dockerProdSamplePath)
	if os.IsNotExist(err) {
		panic("Cannot find docker-compose.sample.yml")
	}
	dockerProdOutputPath := fmt.Sprintf("%s/docker-compose.production.yml", dockerDir)
	genDockerComposeFile(dockerProdSamplePath, dockerProdOutputPath, pubInfos)
}
