// 外部端口监听
package main

type OuterListener struct {
}

func NewOuterListener(cfg *Config) (*OuterListener, error) {
	var listener OuterListener

	return &listener, nil
}
