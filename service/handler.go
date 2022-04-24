package service

import (
	"context"
	. "dotslash/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gofiber/websocket/v2"
)

func (l *Language) Handler(c *websocket.Conn) {

	defer c.Close()

	var (
		msg []byte
		err error
	)

	errorResponse := WsResponse{
		Output:      "",
		Error:       "",
		ServerError: "Server Error\n",
	}

	if _, msg, err = c.ReadMessage(); err != nil {
		log.Println(err)
		return

	}

	body := new(WsBody)
	if err := json.Unmarshal(msg, body); err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}

	directory, err := os.MkdirTemp("", "")
	if err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}
	defer os.RemoveAll(directory)

	if err := os.WriteFile(fmt.Sprintf("%v/%v.%v", directory, l.Filename, l.getExtension()), []byte(body.Code), 0664); err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}

	var dockerClient *client.Client
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var resp container.ContainerCreateCreatedBody
	resp, err = dockerClient.ContainerCreate(timeoutContext, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          false,
		Image:        l.Name,
		WorkingDir:   "/work",
		Cmd:          l.getCommand(body),
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: directory,
				Target: "/work",
			},
		},
	}, nil, nil, "")

	if err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}

	defer func(dockerClient *client.Client) {
		if err := dockerClient.ContainerRemove(timeoutContext, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			log.Println(err)
		}
	}(dockerClient)

	if err = dockerClient.ContainerStart(timeoutContext, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}

	statusCh, errorCh := dockerClient.ContainerWait(timeoutContext, resp.ID, container.WaitConditionNotRunning)

	var containerResp types.HijackedResponse
	containerResp, err = dockerClient.ContainerAttach(timeoutContext, resp.ID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		log.Println(err)
		c.WriteJSON(errorResponse)
		return
	}
	defer containerResp.Close()

	input := make(chan []byte)
	output := make(chan string)
	outputError := make(chan string)
	errChan := make(chan error)
	terminateCh := make(chan bool)

	// Write to docker container
	go func(w io.WriteCloser) {
		for {
			data, ok := <-input
			if !ok {
				w.Close()
				return
			}
			w.Write(data)
		}
	}(containerResp.Conn)

	// Receive from docker container
	go func(reader io.Reader) {
		outputChannelWriter := ChannelWriter{output}
		errorChannelWriter := ChannelWriter{outputError}
		_, err := stdcopy.StdCopy(&outputChannelWriter, &errorChannelWriter, reader)
		if err != nil {
			errChan <- err
		}
	}(containerResp.Reader)

	// Connect STDOUT to websocket
	go func() {
	loop:
		for {
			select {
			case data, ok := <-output:
				if !ok {
					break loop
				}
				if err := c.WriteJSON(WsResponse{
					Output:      data,
					Error:       "",
					ServerError: "",
				}); err != nil {
					errChan <- err
					break loop
				}

			case data, ok := <-outputError:
				if !ok {
					break loop
				}
				if err := c.WriteJSON(WsResponse{
					Output:      "",
					Error:       data,
					ServerError: "",
				}); err != nil {
					errChan <- err
					break loop
				}
			}
		}
	}()

	// Websocket to docker
	go func(c *websocket.Conn) {
		for {
			if _, msg, err := c.ReadMessage(); err != nil {
				errChan <- err
				break
			} else {
				if err := json.Unmarshal(msg, body); err != nil {
					log.Println(err)
					terminateCh <- true
					break
				}
				if body.Interrupt {
					terminateCh <- true
					break
				}
				if body.Input != "" {
					input <- []byte(body.Input)
				}
			}
		}
	}(c)

loop:
	for {
		select {
		case err := <-errorCh:
			log.Println(err)
			c.WriteJSON(errorResponse)
			close(input)
			close(output)
			close(outputError)
			break loop
		case <-statusCh:
			close(input)
			close(output)
			close(outputError)
			break loop
		case err := <-errChan:
			close(input)
			close(output)
			close(outputError)
			log.Println(err)
			c.WriteJSON(errorResponse)
			if err != io.EOF {
				log.Println(err)
			}
			break loop
		case <-terminateCh:
			if err := dockerClient.ContainerStop(timeoutContext, resp.ID, nil); err != nil {
				close(input)
				close(output)
				close(outputError)
				log.Println(err)
			}
			break loop
		}
	}
}
