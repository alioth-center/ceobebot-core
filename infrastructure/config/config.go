package config

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

var (
	configFileContent *bytes.Buffer
)

func init() {
	if loadConfigErr := LoadConfig(); loadConfigErr != nil {
		panic(loadConfigErr)
	}
}

// LoadConfig 加载默认配置文件
func LoadConfig() (err error) {
	configFilePath := "data/config.yaml"
	file, openFileErr := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if openFileErr != nil {
		return fmt.Errorf("open config file error: %w", openFileErr)
	}

	if configFileContentBytes, readFileErr := io.ReadAll(file); readFileErr != nil {
		return fmt.Errorf("read config file error: %w", readFileErr)
	} else {
		configFileContent = bytes.NewBuffer(configFileContentBytes)
	}

	if closeFileErr := file.Close(); closeFileErr != nil {
		return fmt.Errorf("close config file error: %w", closeFileErr)
	}

	return nil
}

// LoadCustomConfig 从默认配置文件中加载自定义配置
//   - receiver: 接收配置的结构体，需要使用指针
//   - key: 配置项的键名
func LoadCustomConfig(receiver any, key string) (err error) {
	buffer := map[string]interface{}{}
	if unmarshalErr := yaml.Unmarshal(configFileContent.Bytes(), &buffer); unmarshalErr != nil {
		return fmt.Errorf("unmarshal content buffer error: %w", unmarshalErr)
	} else if contentBuffer, existKey := buffer[key]; !existKey {
		return NewNoMatchedKeyError(key)
	} else if bytesNewBuffer, marshalErr := yaml.Marshal(contentBuffer); marshalErr != nil {
		return fmt.Errorf("marshal content buffer error: %w", marshalErr)
	} else if unmarshalResultErr := yaml.Unmarshal(bytesNewBuffer, receiver); unmarshalResultErr != nil {
		return fmt.Errorf("unmarshal config content error: %w", unmarshalResultErr)
	} else {
		return nil
	}
}

// LoadCustomConfigWithKeys 从默认配置文件中加载自定义配置，只加载指定的配置项，相当于多次调用 LoadCustomConfig
//   - receiver: 接收配置的结构体，需要使用指针
//   - keys: 配置项的键名，按层级顺序传入
func LoadCustomConfigWithKeys(receiver any, keys ...string) (err error) {
	bytesBuffer := configFileContent.Bytes()
	for _, key := range keys {
		buffer := map[string]interface{}{}
		if unmarshalErr := yaml.Unmarshal(bytesBuffer, &buffer); unmarshalErr != nil {
			return fmt.Errorf("unmarshal content buffer error: %w", unmarshalErr)
		} else if contentBuffer, existKey := buffer[key]; !existKey {
			return NewNoMatchedKeyError(key)
		} else if bytesNewBuffer, marshalErr := yaml.Marshal(contentBuffer); marshalErr != nil {
			return fmt.Errorf("marshal content buffer error: %w", marshalErr)
		} else {
			bytesBuffer = bytesNewBuffer
		}
	}

	if unmarshalErr := yaml.Unmarshal(bytesBuffer, receiver); unmarshalErr != nil {
		return fmt.Errorf("marshal config content error: %w", unmarshalErr)
	}

	return nil
}

// LoadExternalConfig 加载外部配置文件
//   - receiver: 接收配置的结构体，需要使用指针
//   - configFilePath: 配置文件路径，不会检查文件是否存在
func LoadExternalConfig(receiver any, configFilePath string) (err error) {
	file, openFileErr := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if openFileErr != nil {
		return fmt.Errorf("open config file error: %w", openFileErr)
	}

	if configFileContentBytes, readFileErr := io.ReadAll(file); readFileErr != nil {
		return fmt.Errorf("read config file error: %w", readFileErr)
	} else if unmarshalErr := yaml.Unmarshal(configFileContentBytes, receiver); unmarshalErr != nil {
		return fmt.Errorf("unmarshal config file error: %w", unmarshalErr)
	} else if closeFileErr := file.Close(); closeFileErr != nil {
		return fmt.Errorf("close config file error: %w", closeFileErr)
	} else {
		return nil
	}
}

// LoadExternalConfigWithKeys 加载外部配置文件，只加载指定的配置项
//   - receiver: 接收配置的结构体，需要使用指针
//   - configFilePath: 配置文件路径，不会检查文件是否存在
//   - keys: 配置项的键名，按层级顺序传入
func LoadExternalConfigWithKeys(receiver any, configFilePath string, keys ...string) (err error) {
	file, openFileErr := os.OpenFile(configFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if openFileErr != nil {
		return fmt.Errorf("open config file error: %w", openFileErr)
	}

	var bytesBuffer []byte
	if configFileContentBytes, readFileErr := io.ReadAll(file); readFileErr != nil {
		return fmt.Errorf("read config file error: %w", readFileErr)
	} else {
		bytesBuffer = configFileContentBytes
	}

	for _, key := range keys {
		buffer := map[string]interface{}{}
		if unmarshalErr := yaml.Unmarshal(bytesBuffer, &buffer); unmarshalErr != nil {
			return fmt.Errorf("unmarshal content buffer error: %w", unmarshalErr)
		} else if contentBuffer, existKey := buffer[key]; !existKey {
			return NewNoMatchedKeyError(key)
		} else if bytesNewBuffer, marshalErr := yaml.Marshal(contentBuffer); marshalErr != nil {
			return fmt.Errorf("marshal content buffer error: %w", marshalErr)
		} else {
			bytesBuffer = bytesNewBuffer
		}
	}

	if unmarshalErr := yaml.Unmarshal(bytesBuffer, receiver); unmarshalErr != nil {
		return fmt.Errorf("marshal config content error: %w", unmarshalErr)
	}

	return nil
}
