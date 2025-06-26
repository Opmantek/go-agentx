// Copyright 2018 The agentx authors
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package agentx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/pdu"
	"github.com/posteo/go-agentx/value"
)

func TestListHandler(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	session, err := e.client.Session()
	require.NoError(t, err)
	defer session.Close()

	lh := &agentx.ListHandler{}
	i := lh.Add("1.3.6.1.4.1.45995.3.1")
	i.Type = pdu.VariableTypeOctetString
	i.Value = "test"

	//For testing the sort of the OIDS
	var oids = 15
	for j := 0; j < oids; j++ {
		oid := fmt.Sprintf("1.3.6.1.4.1.45995.4.%d", j+1)
		i := lh.Add(oid)
		i.Type = pdu.VariableTypeOctetString
		i.Value = fmt.Sprintf("test%d", j+1)
	}

	// lets register the list handler
	session.Handler = lh

	//lets register over

	baseOID := value.MustParseOID("1.3.6.1.4.1.45995")

	require.NoError(t, session.Register(127, baseOID))
	defer session.Unregister(127, baseOID)

	t.Run("Walk", func(t *testing.T) {
		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.1 = STRING: \"test\"",
			SNMPWalk(t, "1.3.6.1.4.1.45995.3.1"))

		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.2 = No Such Object available on this agent at this OID",
			SNMPWalk(t, "1.3.6.1.4.1.45995.3.2"))
	})

	//	Test to fix issue where walk in order would fail
	//  Error: OID not increasing: .1.3.6.1.4.1.45995.4.15
	t.Run("WalkInOrder", func(t *testing.T) {
		assert.Equal(t,
			".1.3.6.1.4.1.45995.4.1 = STRING: \"test1\"\n.1.3.6.1.4.1.45995.4.2 = STRING: \"test2\"\n.1.3.6.1.4.1.45995.4.3 = STRING: \"test3\"\n.1.3.6.1.4.1.45995.4.4 = STRING: \"test4\"\n.1.3.6.1.4.1.45995.4.5 = STRING: \"test5\"\n.1.3.6.1.4.1.45995.4.6 = STRING: \"test6\"\n.1.3.6.1.4.1.45995.4.7 = STRING: \"test7\"\n.1.3.6.1.4.1.45995.4.8 = STRING: \"test8\"\n.1.3.6.1.4.1.45995.4.9 = STRING: \"test9\"\n.1.3.6.1.4.1.45995.4.10 = STRING: \"test10\"\n.1.3.6.1.4.1.45995.4.11 = STRING: \"test11\"\n.1.3.6.1.4.1.45995.4.12 = STRING: \"test12\"\n.1.3.6.1.4.1.45995.4.13 = STRING: \"test13\"\n.1.3.6.1.4.1.45995.4.14 = STRING: \"test14\"\n.1.3.6.1.4.1.45995.4.15 = STRING: \"test15\"",
			SNMPWalk(t, "1.3.6.1.4.1.45995.4"))

	})

	t.Run("Get", func(t *testing.T) {
		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.1 = STRING: \"test\"",
			SNMPGet(t, []string{"1.3.6.1.4.1.45995.3.1"}))

		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.2 = No Such Object available on this agent at this OID",
			SNMPGet(t, []string{"1.3.6.1.4.1.45995.3.2"}))

		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.1 = STRING: \"test\"\n.1.3.6.1.4.1.45995.4.1 = STRING: \"test1\"",
			SNMPGet(t, []string{"1.3.6.1.4.1.45995.3.1", "1.3.6.1.4.1.45995.4.1"}))
	})

	t.Run("GetNext", func(t *testing.T) {
		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.1 = STRING: \"test\"",
			SNMPGetNext(t, "1.3.6.1.4.1.45995.3.0"))

		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.1 = STRING: \"test\"",
			SNMPGetNext(t, "1.3.6.1.4.1.45995.3"))

	})

	t.Run("GetBulk", func(t *testing.T) {
		assert.Equal(t,
			".1.3.6.1.4.1.45995.3.1 = STRING: \"test\"",
			SNMPGetBulk(t, "1.3.6.1.4.1.45995.3.0", 0, 1))
	})
}
