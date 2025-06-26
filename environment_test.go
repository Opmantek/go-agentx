// Copyright 2018 The agentx authors
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package agentx_test

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/value"
)

type environment struct {
	client   *agentx.Client
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	//Allow to override for use in docker
	snmpCnf := os.Getenv("SNMP_CONF")
	if snmpCnf == "" {
		snmpCnf = "snmpd.conf"
	}
	cmd := exec.Command("snmpd", "-Lo", "-f", "-c", snmpCnf)

	stdout, err := cmd.StdoutPipe()
	require.NoError(tb, err)
	stderr, err := cmd.StderrPipe()
	require.NoError(tb, err)

	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	log.Printf("run: %s", cmd)
	require.NoError(tb, cmd.Start())

	// Wait for snmpd to be ready
	for i := 0; i < 50; i++ { // 5 seconds max
		conn, err := net.DialTimeout("tcp", "127.0.0.1:30705", 100*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	client, err := agentx.Dial("tcp", "127.0.0.1:30705")
	require.NoError(tb, err)
	client.Timeout = 60 * time.Second
	client.ReconnectInterval = 1 * time.Second
	client.NameOID = value.MustParseOID("1.3.6.1.4.1.45995")
	client.Name = "test client"

	return &environment{
		client: client,
		tearDown: func() {
			require.NoError(tb, client.Close())
			require.NoError(tb, cmd.Process.Kill())
			cmd.Wait()
		},
	}
}
