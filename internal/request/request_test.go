//nolint:testpackage
package request

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
)

func TestDoImp(t *testing.T) {
	t.Parallel()

	t.Run("create request error", func(t *testing.T) {
		t.Parallel()

		var (
			req = Request{}
			ctx = context.Background()
		)

		err := req.doImpl(ctx, "/", "", nil)
		require.EqualError(t, err, `create new request: net/http: invalid method "/"`)
	})

	t.Run("doer error", func(t *testing.T) {
		t.Parallel()

		var (
			ctrl        = gomock.NewController(t)
			interceptor = NewMockInterceptor(ctrl)
			doer        = NewMockdoer(ctrl)
			req         = Request{
				interceptor: interceptor,
			}
			ctx = context.Background()
		)

		gomock.InOrder(
			interceptor.EXPECT().
				Before(),
			doer.EXPECT().
				Do(gomock.Any()).
				Return(nil, errors.New("some do error")),
		)

		err := req.doImpl(ctx, "GET", "some_url", doer)
		require.EqualError(t, err, "do http request: some do error")
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ctrl        = gomock.NewController(t)
			interceptor = NewMockInterceptor(ctrl)
			doer        = NewMockdoer(ctrl)
			req         = Request{
				interceptor: interceptor,
			}
			ctx = context.Background()
			rsp = http.Response{
				Body: io.NopCloser(nil),
			}
		)

		gomock.InOrder(
			interceptor.EXPECT().
				Before(),
			doer.EXPECT().
				Do(gomock.Any()).
				Return(&rsp, nil),
			interceptor.EXPECT().
				After(&rsp),
		)

		err := req.doImpl(ctx, "GET", "some_url", doer)
		require.NoError(t, err)
	})
}
