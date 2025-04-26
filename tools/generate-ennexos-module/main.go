package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hansmi/modbus-exporter-modules/internal/ennexos"
	"github.com/hansmi/modbus-exporter-modules/internal/modbusexporter"
)

func transformFiles(w io.Writer, moduleName, config, params string) error {
	cfg, err := loadConfig(config)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	pf, err := ennexos.ReadParameterFile(params)
	if err != nil {
		return fmt.Errorf("loading params: %w", err)
	}

	return transform(w, moduleName, *cfg, *pf)
}

func transform(w io.Writer, moduleName string, cfg config, pf ennexos.ParameterFile) error {
	module := modbusexporter.ModuleDef{
		Name:     moduleName,
		Protocol: modbusexporter.TCPIP,
	}

	parameters, err := preprocess(pf.Parameters)
	if err != nil {
		return err
	}

	module.Metrics, err = parametersToMetrics(cfg, module.Name, parameters)
	if err != nil {
		return err
	}

	config := modbusexporter.Config{
		Modules: []modbusexporter.ModuleDef{module},
	}

	output, err := modbusexporter.MarshalConfig(config)
	if err != nil {
		return err
	}

	if desc := pf.Header.Description; desc != "" {
		fmt.Fprintf(w, "# %s\n\n", desc)
	}

	_, err = w.Write(output)
	return err
}

func main() {
	moduleName := flag.String("module", "",
		"Name of generated modbus-exporter module.")
	configFile := flag.String("config", "",
		"Configuration file in YAML format. Required.")
	paramsFile := flag.String("params", "",
		"Parameter definition file in JSON format. Required.")

	flag.Parse()

	if flag.NArg() > 0 {
		log.Fatal("No arguments expected.")
	}

	if err := transformFiles(os.Stdout, *moduleName, *configFile, *paramsFile); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
