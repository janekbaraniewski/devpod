package ssh

import (
	"context"
	"fmt"
	"io"

	"github.com/loft-sh/devpod/pkg/stdio"
	"github.com/loft-sh/devpod/pkg/tailscale"
	"github.com/loft-sh/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func NewSSHPassClient(user, addr, password string) (*ssh.Client, error) {
	clientConfig := &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	clientConfig.Auth = append(clientConfig.Auth, ssh.Password(password))

	if user != "" {
		clientConfig.User = user
	}

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("dial to %v failed: %w", addr, err)
	}

	return client, nil
}

func NewSSHClient(user, addr string, keyBytes []byte) (*ssh.Client, error) {
	sshConfig, err := ConfigFromKeyBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	if user != "" {
		sshConfig.User = user
	}

	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("dial to %v failed: %w", addr, err)
	}

	return client, nil
}

func StdioClient(reader io.Reader, writer io.WriteCloser, exitOnClose bool) (*ssh.Client, error) {
	return StdioClientFromKeyBytesWithUser(nil, reader, writer, "", exitOnClose)
}

func StdioClientWithUser(reader io.Reader, writer io.WriteCloser, user string, exitOnClose bool) (*ssh.Client, error) {
	return StdioClientFromKeyBytesWithUser(nil, reader, writer, user, exitOnClose)
}

func TailscaleClientWithUser(ctx context.Context, ts tailscale.TSNet, serverAddress, user string, log log.Logger) (*ssh.Client, error) {
	log.Debugf("Connecting to SSH server at %s with user %s", serverAddress, user)

	conn, err := ts.Dial(ctx, "tcp", serverAddress)
	if err != nil {
		log.Errorf("Failed to connect to %s: %v", serverAddress, err)
		return nil, fmt.Errorf("failed to connect to %s: %w", serverAddress, err)
	}

	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{}, // FIXME
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Debugf("Attempting to establish SSH connection with %s as user %s", serverAddress, user)

	sshConn, channels, requests, err := ssh.NewClientConn(conn, serverAddress, clientConfig)
	if err != nil {
		log.Errorf("Failed to establish SSH connection to %s: %v", serverAddress, err)
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	log.Debugf("SSH connection established with %s as user %s", serverAddress, user)
	return ssh.NewClient(sshConn, channels, requests), nil
}

func StdioClientFromKeyBytesWithUser(keyBytes []byte, reader io.Reader, writer io.WriteCloser, user string, exitOnClose bool) (*ssh.Client, error) {
	conn := stdio.NewStdioStream(reader, writer, exitOnClose, 0)
	clientConfig, err := ConfigFromKeyBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	clientConfig.User = user
	c, chans, req, err := ssh.NewClientConn(conn, "stdio", clientConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, req), nil
}

func ConfigFromKeyBytes(keyBytes []byte) (*ssh.ClientConfig, error) {
	clientConfig := &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// key file authentication?
	if len(keyBytes) > 0 {
		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return nil, errors.Wrap(err, "parse private key")
		}

		clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
	}
	return clientConfig, nil
}

func Run(ctx context.Context, client *ssh.Client, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer, envVars map[string]string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	for k, v := range envVars {
		err = sess.Setenv(k, v)
		if err != nil {
			return err
		}
	}

	exit := make(chan struct{})
	defer close(exit)
	go func() {
		select {
		case <-ctx.Done():
			_ = sess.Signal(ssh.SIGINT)
			_ = sess.Close()
		case <-exit:
		}
	}()

	sess.Stdin = stdin
	sess.Stdout = stdout
	sess.Stderr = stderr
	err = sess.Run(command)
	if err != nil {
		return err
	}

	return nil
}
