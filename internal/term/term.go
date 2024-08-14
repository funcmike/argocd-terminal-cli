package term

import (
	"context"
	"errors"
	"fmt"
	"github.com/funcmike/argocd-terminal-cli/internal/argocd"
	"github.com/moby/term"
	"io"
	"log"
	"net/http"
	"os"
)

const cursorGetCode string = "\u001b[6n"

func Run(ctx context.Context, options argocd.TerminalClientOptions, token string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	argoURL, err := argocd.BuildDefaultTerminalURL(options)
	if err != nil {
		return err
	}

	client, err := argocd.NewTerminalClient(ctx, argoURL, argocd.BuildDefaultHeaders(token), http.DefaultClient)
	if err != nil {
		return fmt.Errorf("could not create argocd terminal client: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		defer cancel()
		defer close(done)

		done <- (&terminal{client}).Run(ctx)
	}()

	<-ctx.Done()
	return <-done
}

type terminal struct {
	client *argocd.TerminalClient
}

func (t *terminal) Run(ctx context.Context) error {
	fd := os.Stdout.Fd()
	if !term.IsTerminal(fd) {
		return fmt.Errorf("not a terminal")
	}

	prevState, err := term.SetRawTerminal(fd)
	if err != nil {
		return fmt.Errorf("could not set terminal raw terminal: %w", err)
	}

	defer term.RestoreTerminal(fd, prevState)

	stdin, stdout, _ := term.StdStreams()
	defer stdin.Close()

	startWinsize := t.getWinsize(fd)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()

		// stdin is a blocking goroutine and take wait forever
		// so just log error
		if err := t.processStdin(ctx, stdin, fd, startWinsize); err != nil {
			log.Printf("could not process stdin: %v", err)
		}
	}()

	return t.processStdout(ctx, stdout, startWinsize)
}

func (t *terminal) processStdin(ctx context.Context, stdin io.ReadCloser, fd uintptr, startWinsize *term.Winsize) error {
	bytes := make([]byte, 1024)
	prevWinsize := startWinsize
	done := ctx.Done()

	for {
		select {
		case <-done:
			return nil
		default:
		}

		n, err := stdin.Read(bytes)

		if n == 0 && err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}

		currWinsize := t.getWinsize(fd)
		windowChanged := prevWinsize.Width != currWinsize.Width || prevWinsize.Height != currWinsize.Height

		if windowChanged {
			prevWinsize = currWinsize

			if err := t.sendResize(ctx, currWinsize); err != nil {
				return fmt.Errorf("error sending resize: %w", err)
			}
		}

		if err := t.sendStdin(ctx, currWinsize, string(bytes[:n])); err != nil {
			return fmt.Errorf("error sending stdin: %w", err)
		}
	}
}

func (t *terminal) processStdout(ctx context.Context, stdout io.Writer, startWinsize *term.Winsize) error {
	var first bool

	for {
		op, err := t.client.Recv(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("could not receive operation: %w", err)
		}

		if op.Operation != argocd.OpStdout {
			return fmt.Errorf("unexpected operation: %s", op.Operation)
		}

		if first {
			if err := t.sendResizeWithCursorPosition(ctx, startWinsize, op); err != nil {
				return fmt.Errorf("error sending first: %w", err)
			}
			first = false
		}

		if _, err := fmt.Fprintf(stdout, op.Data); err != nil {
			return fmt.Errorf("could not write to stdout: %w", err)
		}
	}
}

func (t *terminal) getWinsize(fd uintptr) *term.Winsize {
	size, err := term.GetWinsize(fd)
	if err != nil {
		// should always return size
		panic(err)
	}
	return size
}

func (t *terminal) sendStdin(ctx context.Context, winsize *term.Winsize, data string) error {
	return t.client.Send(ctx, argocd.Operation{
		Operation: argocd.OpStdin,
		Data:      data,
		Rows:      int(winsize.Height),
		Cols:      int(winsize.Width),
	})
}

func (t *terminal) sendResize(ctx context.Context, winsize *term.Winsize) error {
	return t.client.Send(ctx, argocd.Operation{
		Operation: argocd.OpResize,
		Rows:      int(winsize.Height),
		Cols:      int(winsize.Width),
	})
}

func (t *terminal) sendResizeWithCursorPosition(ctx context.Context, winsize *term.Winsize, op argocd.Operation) error {
	if err := t.sendResize(ctx, winsize); err != nil {
		return fmt.Errorf("error sending resize: %w", err)
	}

	if err := t.sendStdin(ctx, winsize, fmt.Sprintf("\u001b[1;%dR", len(op.Data)-len(cursorGetCode))); err != nil {
		return fmt.Errorf("error sending cursor position: %w", err)
	}
	return nil
}
