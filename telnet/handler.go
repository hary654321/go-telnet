/*
 * @Description:
 * @Version: 2.0
 * @Autor: ABing
 * @Date: 2024-06-28 11:07:36
 * @LastEditors: lhl
 * @LastEditTime: 2024-06-28 11:51:35
 */
package telnet

import (
	"context"
	"io"
)

// A Handler serves a TELNET (or TELNETS) connection.
//
// Writing data to the Writer passed as an argument to the ServeTELNET method
// will send data to the TELNET (or TELNETS) client.
//
// Reading data from the Reader passed as an argument to the ServeTELNET method
// will receive data from the TELNET client.
//
// The Writer's Write method sends "escaped" TELNET (and TELNETS) data.
//
// The Reader's Read method "un-escapes" TELNET (and TELNETS) data, and filters
// out TELNET (and TELNETS) command sequences.
type Handler interface {
	ServeTELNET(ctx context.Context, w io.Writer, r io.Reader) error
}
