package goka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/lovoo/goka/multierr"
)

func TestMultiErr_Unwrap(t *testing.T) {

	errs := new(multierr.Errors)

	errs.Collect(fmt.Errorf("other error"))
	errs.Collect(newErrProcessing(0, fmt.Errorf("hello world")))
	errs.Collect(fmt.Errorf("last error"))

	var cbErr *errProcessing
	if errors.As(errs.NilOrError(), &cbErr) {
		log.Printf("is processing error: %v", cbErr)
	} else {
		log.Printf("no processin error")
	}
}

func TestHashuMultierr_Unwrap(t *testing.T) {

	errs := multierror.Append(fmt.Errorf("other error"),
		fmt.Errorf("last error"),
		fmt.Errorf("some wrapping error %w", fmt.Errorf("some error %w", newErrProcessing(1, fmt.Errorf("hello world")))),
	)

	var cbErr *errProcessing
	if errors.As(errs, &cbErr) {
		log.Printf("is processing error: %v", cbErr)
	} else {
		log.Printf("no processing error")
	}
}

func TestMultiError_Group(t *testing.T) {

	g, ctx := multierr.NewErrGroup(context.Background())

	g.Go(func() error {
		select {
		case <-time.After(10 * time.Millisecond):
			log.Printf("success")
			return nil
		case <-ctx.Done():
			return fmt.Errorf("A: context done")
		}
	})
	g.Go(func() error {
		select {
		case <-time.After(100 * time.Millisecond):
			return fmt.Errorf("A: context not done")
		case <-ctx.Done():
			return fmt.Errorf("A: context done")
		}
	})
	g.Go(func() error {
		select {
		case <-time.After(50 * time.Millisecond):
			return fmt.Errorf("B: context not done")
		case <-ctx.Done():
			return fmt.Errorf("B: context done")
		}
	})

	log.Printf("%v", g.Wait().ErrorOrNil())

}
